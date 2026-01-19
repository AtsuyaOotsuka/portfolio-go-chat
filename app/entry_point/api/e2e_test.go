package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/api_mock"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var baseURL string
var mongo *usecase.Mongo
var mongoHelper *funcs.TestMongoStruct

func mockApi() *echo.Echo {
	apiMock := echo.New()
	apiMock.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Printf(
				"[MOCK API] %s %s",
				c.Request().Method,
				c.Request().URL.Path,
			)
			return next(c)
		}
	})

	group := apiMock.Group("/server_api")
	group.POST("/user/profile", api_mock.AuthUserGetProfile)
	return apiMock
}
func waitForServer(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", addr)
}

func TestMain(m *testing.M) {
	var err error

	mongo, err = SetupMongo()
	if err != nil {
		panic(err)
	}

	mongoHelper = funcs.SetUpMongoTestDatabase()
	defer mongoHelper.Disconnect()

	baseURL = "http://127.0.0.1:8880"

	redis, err := SetupRedis()
	if err != nil {
		panic(err)
	}

	apiMock := mockApi()

	go func() {
		if err := apiMock.Start("127.0.0.1:8881"); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	app := SetupRouter(mongo, redis)
	defer app.Shutdown()

	testServer := &http.Server{
		Addr: "127.0.0.1:8880",
	}

	go func() {
		if err := app.Echo.StartServer(testServer); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	defer testServer.Close()

	if err := waitForServer("127.0.0.1:8881", 5*time.Second); err != nil {
		panic(err)
	}
	if err := waitForServer("127.0.0.1:8880", 5*time.Second); err != nil {
		panic(err)
	}

	exitCode := m.Run()
	if err := testServer.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

func createCsrf() string {
	csrf_token := os.Getenv("CSRF_TOKEN")
	nonce := funcs.GenerateCSRFCookieToken(
		csrf_token,
		time.Now().Add(1*time.Hour).Unix(),
	)
	return nonce
}

func createJwt(
	uuid string,
	Email string,
	Exp time.Time,
) string {
	jwt_token := os.Getenv("JWT_SECRET_KEY")

	jwtSvc := atylabjwt.NewJwtSvc()
	jwtConfig := &atylabjwt.JwtConfig{
		Key:   []byte(jwt_token),
		Uuid:  uuid,
		Email: Email,
		Exp:   Exp,
	}
	jwt, err := jwtSvc.CreateJwt(jwtConfig)
	if err != nil {
		panic(err)
	}

	return jwt
}

func request(method string, url string, jwt string, body io.Reader, t *testing.T) (*http.Response, func() error) {
	csrf := createCsrf()

	client := &http.Client{}
	requestUrl := baseURL + url
	fmt.Println("Request URL:", requestUrl)
	req, err := http.NewRequest(method, requestUrl, body)
	if method != "GET" && err == nil {
		req.Header.Set("Content-Type", "application/json")
	}
	assert.NoError(t, err)
	if method != "GET" {
		req.Header.Set("X-CSRF-Token", csrf)
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	resp, err := client.Do(req)
	assert.NoError(t, err)

	return resp, resp.Body.Close
}

func TestHealth(t *testing.T) {
	jwt := createJwt(
		"test-uuid",
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	resp, close := request("GET", "/healthcheck", jwt, nil, t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
}

func TestRoomList(t *testing.T) {
	mongoHelper.MongoCleanUp()

	variations := []model.Room{
		{Name: "PrivateRoom_Owner", OwnerID: "test-uuid", IsPrivate: true, Members: []string{"test-uuid", "99999"}},
		{Name: "PrivateRoom_Member", OwnerID: "99999", IsPrivate: true, Members: []string{"99999", "test-uuid"}},
		{Name: "PrivateRoom_None", OwnerID: "88888", IsPrivate: true, Members: []string{"88888"}},
		{Name: "PublicRoom_None", OwnerID: "77777", IsPrivate: false, Members: []string{"77777"}},
		{Name: "PublicRoom_Joined", OwnerID: "66666", IsPrivate: false, Members: []string{"66666", "test-uuid"}},
	}

	for i := range variations {
		InsertId, err := mongoHelper.Insert(
			model.RoomCollectionName, variations[i],
		)
		assert.NoError(t, err)
		variations[i].ID, _ = primitive.ObjectIDFromHex(InsertId)
	}

	uuid := "test-uuid"

	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)

	expected := map[string]map[string]any{
		"success all": {
			"target":    "all",
			"views":     []string{"PublicRoom_None", "PrivateRoom_Owner", "PrivateRoom_Member", "PublicRoom_Joined"},
			"not_views": []string{"PrivateRoom_None"},
		},
		"success joined": {
			"target":    "joined",
			"views":     []string{"PrivateRoom_Owner", "PrivateRoom_Member", "PublicRoom_Joined"},
			"not_views": []string{"PublicRoom_None", "PrivateRoom_None"},
		},
		"success default(all)": {
			"target":    "",
			"views":     []string{"PublicRoom_None", "PrivateRoom_Owner", "PrivateRoom_Member", "PublicRoom_Joined"},
			"not_views": []string{"PrivateRoom_None"},
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {
			resp, close := request("GET", "/room/list?target="+expect["target"].(string), jwt, nil, t)
			defer close()

			assert.Equal(t, 200, resp.StatusCode)
			bodyBytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			bodyString := string(bodyBytes)

			result := map[string][]interface{}{}
			err = json.Unmarshal(bodyBytes, &result)
			assert.NoError(t, err)
			count := len(expect["views"].([]string))
			assert.Len(t, result["rooms"], count)

			// 表示されるルームを確認
			for _, v := range expect["views"].([]string) {
				assert.Contains(t, bodyString, v)
			}
			// 表示されないルームを確認
			for _, v := range expect["not_views"].([]string) {
				assert.NotContains(t, bodyString, v)
			}
		})
	}
}

func TestCreateRoom(t *testing.T) {
	mongoHelper.MongoCleanUp()

	uuid := "test-uuid"

	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)

	requestBody := `{
			"name": "New Public Room",
			"is_private": false
		}`
	resp, close := request("POST", "/room/create", jwt, io.NopCloser(io.Reader(strings.NewReader(requestBody))), t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string]string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result["room_id"])

	exists, err := mongoHelper.ExistContents(model.RoomCollectionName, bson.M{"_id": func() primitive.ObjectID {
		id, _ := primitive.ObjectIDFromHex(result["room_id"])
		return id
	}()})
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRoomJoin(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Joinable Room",
		OwnerID:   "owner-uuid",
		IsPrivate: false,
		Members:   []string{"owner-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	uuid := "test-uuid"

	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	requestBody := fmt.Sprintf(`{
			"room_id": "%s"
		}`, roomID)
	resp, close := request("POST", "/room/"+roomID+"/join", jwt, io.NopCloser(io.Reader(strings.NewReader(requestBody))), t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string]string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, "Joined room successfully", result["message"])

	// メンバーに追加されていることを確認
	var updatedRoom model.Room
	singleResult, err := mongoHelper.FindOneContents(
		model.RoomCollectionName,
		roomID,
	)
	assert.NoError(t, err)
	err = singleResult.Decode(&updatedRoom)
	assert.NoError(t, err)

	assert.Contains(t, updatedRoom.Members, uuid)
}

func TestRoomMembers(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Member List Room",
		OwnerID:   "owner-uuid",
		IsPrivate: false,
		Members:   []string{"owner-uuid", "test-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	uuid := "test-uuid"
	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	resp, close := request("GET", "/room/"+roomID+"/members", jwt, nil, t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
}

func TestRoomLeave(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Leave Room",
		OwnerID:   "owner-uuid",
		IsPrivate: false,
		Members:   []string{"owner-uuid", "test-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	uuid := "test-uuid"
	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	resp, close := request("POST", "/room/"+roomID+"/leave", jwt, nil, t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string]string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, "left room", result["message"])

	// メンバーから削除されていることを確認
	var updatedRoom model.Room
	singleResult, err := mongoHelper.FindOneContents(
		model.RoomCollectionName,
		roomID,
	)
	assert.NoError(t, err)
	err = singleResult.Decode(&updatedRoom)
	assert.NoError(t, err)

	assert.NotContains(t, updatedRoom.Members, "test-uuid")
}

func TestRoomDelete(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Deletable Room",
		OwnerID:   "test-uuid",
		IsPrivate: false,
		Members:   []string{"test-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	uuid := "test-uuid"
	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	resp, close := request("DELETE", "/room/"+roomID+"/admin/delete", jwt, nil, t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string]string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, "room deleted", result["message"])

	// ルームが削除されていることを確認
	exists, err := mongoHelper.ExistContents(model.RoomCollectionName, bson.M{"_id": func() primitive.ObjectID {
		id, _ := primitive.ObjectIDFromHex(roomID)
		return id
	}()})
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestRoomAddMember(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Add Member Room",
		OwnerID:   "test-uuid",
		IsPrivate: false,
		Members:   []string{"test-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	uuid := "test-uuid"
	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	requestBody := fmt.Sprintf(`{
			"member_id": "%s"
		}`, uuid)
	resp, close := request("POST", "/room/"+roomID+"/admin/add_member", jwt, io.NopCloser(io.Reader(strings.NewReader(requestBody))), t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string]string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, "member added", result["message"])

	// メンバーに追加されていることを確認
	var updatedRoom model.Room
	singleResult, err := mongoHelper.FindOneContents(
		model.RoomCollectionName,
		roomID,
	)
	assert.NoError(t, err)
	err = singleResult.Decode(&updatedRoom)
	assert.NoError(t, err)

	assert.Contains(t, updatedRoom.Members, ""+uuid)
}

func TestRoomRemoveMember(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Remove Member Room",
		OwnerID:   "test-uuid",
		IsPrivate: false,
		Members:   []string{"test-uuid", "test-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	uuid := "test-uuid"
	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	requestBody := fmt.Sprintf(`{
			"member_id": "%s"
		}`, uuid)
	resp, close := request("DELETE", "/room/"+roomID+"/admin/remove_member", jwt, io.NopCloser(io.Reader(strings.NewReader(requestBody))), t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string]string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, "member removed", result["message"])

	// メンバーから削除されていることを確認
	var updatedRoom model.Room
	singleResult, err := mongoHelper.FindOneContents(
		model.RoomCollectionName,
		roomID,
	)
	assert.NoError(t, err)
	err = singleResult.Decode(&updatedRoom)
	assert.NoError(t, err)

	assert.NotContains(t, updatedRoom.Members, uuid)
}

func TestMessageList(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Message List Room",
		OwnerID:   "test-uuid",
		IsPrivate: false,
		Members:   []string{"test-uuid", "test-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	messages := []model.Message{
		{RoomID: roomID, Sender: "test-uuid", Message: "Hello!", CreatedAt: time.Now()},
		{RoomID: roomID, Sender: "test-uuid", Message: "Hi there!", CreatedAt: time.Now()},
	}
	for _, msg := range messages {
		_, err := mongoHelper.Insert(
			model.MessageCollectionName,
			msg,
		)
		assert.NoError(t, err)
	}
	uuid := "test-uuid"
	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	resp, close := request("GET", "/message/"+roomID+"/list", jwt, nil, t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string][]map[string]any{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Len(t, result["messages"], 2)
	assert.Equal(t, "Hello!", result["messages"][0]["Message"])
	assert.Equal(t, "Hi there!", result["messages"][1]["Message"])
}

func TestMessageSend(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Message Send Room",
		OwnerID:   "test-uuid",
		IsPrivate: false,
		Members:   []string{"test-uuid", "test-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	uuid := "test-uuid"
	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	requestBody := `{
			"message": "This is a test message."
		}`
	resp, close := request("POST", "/message/"+roomID+"/send", jwt, io.NopCloser(io.Reader(strings.NewReader(requestBody))), t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string]string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)

	// メッセージが保存されていることを確認
	exists, err := mongoHelper.ExistContents(model.MessageCollectionName, bson.M{"message": "This is a test message."})
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestMessageRead(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Message Read Room",
		OwnerID:   "test-uuid",
		IsPrivate: false,
		Members:   []string{"test-uuid", "sender-test-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	message := model.Message{
		RoomID:        roomID,
		Sender:        "sender-test-uuid",
		Message:       "Please read this message.",
		CreatedAt:     time.Now(),
		IsReadUserIds: []string{},
	}
	messageID, err := mongoHelper.Insert(
		model.MessageCollectionName,
		message,
	)
	assert.NoError(t, err)

	uuid := "test-uuid"
	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	requestBody := fmt.Sprintf(`{
			"message_ids": ["%s"]
		}`, messageID)
	resp, close := request("POST", "/message/"+roomID+"/read", jwt, io.NopCloser(io.Reader(strings.NewReader(requestBody))), t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string]string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)

	// メッセージのIsReadUserIdsにユーザーIDが追加されていることを確認
	var updatedMessage model.Message
	singleResult, err := mongoHelper.FindOneContents(
		model.MessageCollectionName,
		messageID,
	)
	assert.NoError(t, err)
	err = singleResult.Decode(&updatedMessage)
	assert.NoError(t, err)

	assert.Contains(t, updatedMessage.IsReadUserIds, ""+uuid)
}

func TestMessageDelete(t *testing.T) {
	var err error
	mongoHelper.MongoCleanUp()

	room := model.Room{
		Name:      "Message Delete Room",
		OwnerID:   "test-uuid",
		IsPrivate: false,
		Members:   []string{"test-uuid", "sender-test-uuid"},
		CreatedAt: time.Now(),
	}
	roomID, err := mongoHelper.Insert(
		model.RoomCollectionName,
		room,
	)
	assert.NoError(t, err)

	message := model.Message{
		RoomID:    roomID,
		Sender:    "sender-test-uuid",
		Message:   "This message will be deleted.",
		CreatedAt: time.Now(),
	}
	messageID, err := mongoHelper.Insert(
		model.MessageCollectionName,
		message,
	)
	assert.NoError(t, err)

	uuid := "test-uuid"
	jwt := createJwt(
		uuid,
		"test@example.com",
		time.Now().Add(1*time.Hour),
	)
	requestBody := fmt.Sprintf(`{
			"message_id": "%s"
		}`, messageID)
	resp, close := request("DELETE", "/message/"+roomID+"/delete", jwt, io.NopCloser(io.Reader(strings.NewReader(requestBody))), t)
	defer close()

	assert.Equal(t, 200, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	result := map[string]string{}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)

	// メッセージが削除されていることを確認
	exists, err := mongoHelper.ExistContents(model.MessageCollectionName, bson.M{"_id": func() primitive.ObjectID {
		id, _ := primitive.ObjectIDFromHex(messageID)
		return id
	}()})
	assert.NoError(t, err)
	assert.False(t, exists)
}
