package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/dto"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/svc_mock/mongo_svc_mock"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRoomList(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"target":             "all",
			"expect_target":      "all",
			"status":             200,
			"GetRoomListCalled":  1,
			"GetRoomListSuccess": true,
			"success":            true,
		},
		"success with target joined": {
			"target":             "joined",
			"expect_target":      "joined",
			"status":             200,
			"GetRoomListCalled":  1,
			"GetRoomListSuccess": true,
			"success":            true,
		},
		"success with target none": {
			"target":             "",
			"expect_target":      "all",
			"status":             200,
			"GetRoomListCalled":  1,
			"GetRoomListSuccess": true,
			"success":            true,
		},
		"failure to get room list": {
			"target":             "all",
			"expect_target":      "all",
			"status":             500,
			"GetRoomListCalled":  1,
			"GetRoomListSuccess": false,
			"success":            false,
		},
	}

	var err error

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/room/list?target="+expect["target"].(string), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("uuid", "test-uuid-1234")

			dto := dto.NewRoomDtoStruct()
			svcMock := new(mongo_svc_mock.RoomSvcMock)

			returnData := []model.Room{
				{
					ID:        primitive.NewObjectID(),
					Name:      "Test Room",
					OwnerID:   "owner-uuid-5678",
					IsPrivate: false,
					Members:   []string{"test-uuid-1234", "member-uuid-91011"},
					CreatedAt: time.Now(),
				},
				{
					ID:        primitive.NewObjectID(),
					Name:      "Private Room",
					OwnerID:   "owner-uuid-5678",
					IsPrivate: true,
					Members:   []string{"member-uuid-91011"},
					CreatedAt: time.Now(),
				},
			}
			if !expect["GetRoomListSuccess"].(bool) {
				returnData = []model.Room{}
			}

			var returnErr error = nil
			if !expect["GetRoomListSuccess"].(bool) {
				returnErr = fmt.Errorf("GetRoomList error")
			}

			svcMock.On("GetRoomList", "test-uuid-1234", expect["expect_target"].(string), mock.Anything).Return(returnData, returnErr).Times(expect["GetRoomListCalled"].(int))

			handler := NewRoomHandler(svcMock, dto)
			err = handler.List(c)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)
			svcMock.AssertExpectations(t)

			if expect["status"].(int) != http.StatusOK {
				return
			}

			result := map[string][]interface{}{}
			err = json.Unmarshal(rec.Body.Bytes(), &result)
			assert.NoError(t, err)

			assert.Len(t, result["rooms"], 2)

			assert.Equal(t, "Test Room", result["rooms"][0].(map[string]interface{})["Name"])
			assert.Equal(t, false, result["rooms"][0].(map[string]interface{})["IsPrivate"])
			assert.Equal(t, true, result["rooms"][0].(map[string]interface{})["IsMember"])

			assert.Equal(t, "Private Room", result["rooms"][1].(map[string]interface{})["Name"])
			assert.Equal(t, true, result["rooms"][1].(map[string]interface{})["IsPrivate"])
			assert.Equal(t, false, result["rooms"][1].(map[string]interface{})["IsMember"])
		})
	}
}

func TestRoomCreate(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status": 200,
			"body": map[string]interface{}{
				"name":       "New Room",
				"is_private": false,
			},
			"success":           true,
			"createRoomCalled":  1,
			"createRoomSuccess": true,
		},
		"validation error (missing name)": {
			"status": 400,
			"body": map[string]interface{}{
				"is_private": false,
			},
			"success":           false,
			"createRoomCalled":  0,
			"createRoomSuccess": false,
		},
		"failure to create room": {
			"status": 500,
			"body": map[string]interface{}{
				"name":       "New Room",
				"is_private": true,
			},
			"success":           false,
			"createRoomCalled":  1,
			"createRoomSuccess": false,
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &usecase.CustomValidator{Validator: validator.New()}

			body := expect["body"].(map[string]interface{})
			jsonBody, _ := json.Marshal(body)
			reqBody := strings.NewReader(string(jsonBody))

			req := httptest.NewRequest(http.MethodPost, "/room/create", reqBody)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("uuid", "test-uuid-1234")

			dto := dto.NewRoomDtoStruct()

			var roomId string = "new-room-id-5678"
			var returnErr error = nil
			if !expect["createRoomSuccess"].(bool) {
				roomId = ""
				returnErr = fmt.Errorf("CreateRoom error")
			}

			svcMock := new(mongo_svc_mock.RoomSvcMock)

			if expect["createRoomCalled"].(int) != 0 {
				svcMock.On("CreateRoom", mock.AnythingOfType("model.Room"), mock.Anything).Return(roomId, returnErr).Times(expect["createRoomCalled"].(int))
			}

			handler := NewRoomHandler(svcMock, dto)
			err := handler.Create(c)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["createRoomCalled"].(int) == 0 {
				svcMock.AssertNotCalled(t, "CreateRoom")
			} else {
				svcMock.AssertExpectations(t)
			}

			if expect["status"].(int) != http.StatusOK {
				return
			}

			result := map[string]interface{}{}
			err = json.Unmarshal(rec.Body.Bytes(), &result)
			assert.NoError(t, err)

			assert.Equal(t, "new-room-id-5678", result["room_id"])
			assert.Equal(t, "New Room", result["room_name"])
			assert.NotEmpty(t, result["created_at"])
		})
	}
}

func TestRoomJoin(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status":          200,
			"body":            map[string]interface{}{"room_id": "existing-room-id-1234"},
			"JoinRoomCalled":  1,
			"JoinRoomSuccess": true,
			"is_member":       false,
		},
		"validation error (missing room_id)": {
			"status":          400,
			"body":            map[string]interface{}{},
			"JoinRoomCalled":  0,
			"JoinRoomSuccess": false,
			"is_member":       false,
		},
		"already a member": {
			"status":          400,
			"body":            map[string]interface{}{"room_id": "existing-room-id-1234"},
			"JoinRoomCalled":  0,
			"JoinRoomSuccess": false,
			"is_member":       true,
		},
		"failure to join room": {
			"status":          500,
			"body":            map[string]interface{}{"room_id": "existing-room-id-1234"},
			"JoinRoomCalled":  1,
			"JoinRoomSuccess": false,
			"is_member":       false,
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {

			e := echo.New()
			e.Validator = &usecase.CustomValidator{Validator: validator.New()}

			body := expect["body"].(map[string]interface{})

			jsonBody, _ := json.Marshal(body)
			reqBody := strings.NewReader(string(jsonBody))

			req := httptest.NewRequest(http.MethodPost, "/room/join", reqBody)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("uuid", "test-uuid-1234")

			room := model.Room{
				ID:      primitive.NewObjectID(),
				Name:    "Test Room",
				OwnerID: "owner-uuid-5678",
				Members: []string{"test-uuid-1234", "another-uuid-91011"},
			}
			c.Set("room_model", room)
			c.Set("is_member", expect["is_member"].(bool))

			dto := dto.NewRoomDtoStruct()
			svcMock := new(mongo_svc_mock.RoomSvcMock)
			if expect["JoinRoomCalled"].(int) != 0 {
				var returnErr error = nil
				if !expect["JoinRoomSuccess"].(bool) {
					returnErr = fmt.Errorf("JoinRoom error")
				}
				svcMock.On("JoinRoom", "existing-room-id-1234", "test-uuid-1234", mock.Anything).Return(returnErr).Times(expect["JoinRoomCalled"].(int))
			}

			handler := NewRoomHandler(svcMock, dto)
			err := handler.Join(c)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["JoinRoomCalled"].(int) != 0 {
				svcMock.AssertExpectations(t)
			} else {
				svcMock.AssertNotCalled(t, "JoinRoom")
			}
		})
	}

}
