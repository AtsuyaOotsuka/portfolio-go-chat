package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
	"github.com/stretchr/testify/assert"
)

var baseURL string
var mongo *usecase.Mongo
var mongoHelper *funcs.TestMongoStruct

func TestMain(m *testing.M) {
	var err error

	mongo, err = SetupMongo()
	if err != nil {
		panic(err)
	}

	mongoHelper = funcs.SetUpMongoTestDatabase()
	defer mongoHelper.Disconnect()

	baseURL = "http://localhost:8880"

	app := SetupRouter(mongo)
	defer app.Shutdown()

	testServer := &http.Server{
		Addr: ":8880",
	}

	go func() {
		if err := app.Echo.StartServer(testServer); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	defer testServer.Close()

	time.Sleep(200 * time.Millisecond)

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
		{Name: "PrivateRoom_Owner", OwnerID: "usertest-uuid", IsPrivate: true, Members: []string{"usertest-uuid", "99999"}},
		{Name: "PrivateRoom_Member", OwnerID: "99999", IsPrivate: true, Members: []string{"99999", "usertest-uuid"}},
		{Name: "PrivateRoom_None", OwnerID: "88888", IsPrivate: true, Members: []string{"88888"}},
		{Name: "PublicRoom_None", OwnerID: "77777", IsPrivate: false, Members: []string{"77777"}},
		{Name: "PublicRoom_Joined", OwnerID: "66666", IsPrivate: false, Members: []string{"66666", "usertest-uuid"}},
	}
	_, err := mongoHelper.InsertRooms(variations)
	assert.NoError(t, err)

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
