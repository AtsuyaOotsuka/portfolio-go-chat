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
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/svc_mock"
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
			mongoSvcMock := new(mongo_svc_mock.RoomSvcMock)
			roomSvcMock := new(svc_mock.RoomSvcMock)

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

			mongoSvcMock.On("GetRoomList", "test-uuid-1234", expect["expect_target"].(string), mock.Anything).Return(returnData, returnErr).Times(expect["GetRoomListCalled"].(int))

			handler := NewRoomHandler(mongoSvcMock, roomSvcMock, dto)
			err = handler.List(c)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)
			mongoSvcMock.AssertExpectations(t)

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

			mongoSvcMock := new(mongo_svc_mock.RoomSvcMock)
			roomSvcMock := new(svc_mock.RoomSvcMock)

			if expect["createRoomCalled"].(int) != 0 {
				mongoSvcMock.On("CreateRoom", mock.AnythingOfType("model.Room"), mock.Anything).Return(roomId, returnErr).Times(expect["createRoomCalled"].(int))
			}

			handler := NewRoomHandler(mongoSvcMock, roomSvcMock, dto)
			err := handler.Create(c)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["createRoomCalled"].(int) == 0 {
				mongoSvcMock.AssertNotCalled(t, "CreateRoom")
			} else {
				mongoSvcMock.AssertExpectations(t)
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
			mongoSvcMock := new(mongo_svc_mock.RoomSvcMock)
			roomSvcMock := new(svc_mock.RoomSvcMock)

			if expect["JoinRoomCalled"].(int) != 0 {
				var returnErr error = nil
				if !expect["JoinRoomSuccess"].(bool) {
					returnErr = fmt.Errorf("JoinRoom error")
				}
				mongoSvcMock.On("JoinRoom", "existing-room-id-1234", "test-uuid-1234", mock.Anything).Return(returnErr).Times(expect["JoinRoomCalled"].(int))
			}

			handler := NewRoomHandler(mongoSvcMock, roomSvcMock, dto)
			err := handler.Join(c)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["JoinRoomCalled"].(int) != 0 {
				mongoSvcMock.AssertExpectations(t)
			} else {
				mongoSvcMock.AssertNotCalled(t, "JoinRoom")
			}
		})
	}
}

func TestRoomMembers(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status":   200,
			"error":    nil,
			"IsMember": true,
		},
		"validation error (not a member)": {
			"status":   400,
			"error":    nil,
			"IsMember": false,
		},
		"failure to get member infos": {
			"status":   500,
			"error":    fmt.Errorf("GetMemberInfos error"),
			"IsMember": true,
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {

			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/room/:room_id/members", nil)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("uuid", "test-uuid-1234")
			c.Set("is_member", expect["IsMember"].(bool))

			room := model.Room{
				ID:      primitive.NewObjectID(),
				Name:    "Test Room",
				OwnerID: "owner-uuid-5678",
				Members: []string{"test-uuid-1234", "another-uuid-91011"},
			}
			c.Set("room_model", room)

			c.SetParamNames("room_id")
			c.SetParamValues("test-room-id")

			dto := dto.NewRoomDtoStruct()
			mongoSvcMock := new(mongo_svc_mock.RoomSvcMock)
			roomSvcMock := new(svc_mock.RoomSvcMock)
			roomSvcMock.On("GetMemberInfos", room, mock.Anything).Return(
				[]model.RoomMember{
					{
						Uuid:  "test-uuid",
						Name:  "Test User",
						Email: "test@example.com",
					},
					{
						Uuid:  "owner-uuid",
						Name:  "Owner User",
						Email: "owner@example.com",
					},
				}, expect["error"]).Once()

			handler := NewRoomHandler(mongoSvcMock, roomSvcMock, dto)
			err := handler.Members(c)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["status"].(int) != 200 {
				return
			}

			result := map[string][]interface{}{}
			err = json.Unmarshal(rec.Body.Bytes(), &result)
			assert.NoError(t, err)

			assert.Len(t, result["members"], 2)
			assert.Equal(t, "Test User", result["members"][0].(map[string]interface{})["name"])
			assert.Equal(t, "Owner User", result["members"][1].(map[string]interface{})["name"])

			roomSvcMock.AssertExpectations(t)
		})
	}
}

func TestRoomLeave(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status":           200,
			"IsMember":         true,
			"IsAdmin":          false,
			"LeaveRoomCalled":  1,
			"LeaveRoomSuccess": true,
		},
		"validation error (not a member)": {
			"status":           400,
			"IsMember":         false,
			"IsAdmin":          false,
			"LeaveRoomCalled":  0,
			"LeaveRoomSuccess": false,
		},
		"validation error (is admin)": {
			"status":           400,
			"IsMember":         true,
			"IsAdmin":          true,
			"LeaveRoomCalled":  0,
			"LeaveRoomSuccess": false,
		},
		"failure to leave room": {
			"status":           500,
			"IsMember":         true,
			"IsAdmin":          false,
			"LeaveRoomCalled":  1,
			"LeaveRoomSuccess": false,
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {

			e := echo.New()

			req := httptest.NewRequest(http.MethodPost, "/room/:room_id/leave", nil)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("uuid", "test-uuid-1234")
			c.Set("is_member", expect["IsMember"].(bool))
			c.Set("is_admin", expect["IsAdmin"].(bool))

			c.SetParamNames("room_id")
			c.SetParamValues("test-room-id")

			dto := dto.NewRoomDtoStruct()
			mongoSvcMock := new(mongo_svc_mock.RoomSvcMock)
			roomSvcMock := new(svc_mock.RoomSvcMock)

			var returnErr error = nil
			if !expect["LeaveRoomSuccess"].(bool) {
				returnErr = fmt.Errorf("LeaveRoom error")
			}
			mongoSvcMock.On("LeaveRoom", "test-room-id", "test-uuid-1234", mock.Anything).Return(returnErr).Times(expect["LeaveRoomCalled"].(int))

			handler := NewRoomHandler(mongoSvcMock, roomSvcMock, dto)
			err := handler.Leave(c)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["LeaveRoomCalled"].(int) != 0 {
				mongoSvcMock.AssertExpectations(t)
			} else {
				mongoSvcMock.AssertNotCalled(t, "LeaveRoom")
			}
		})
	}
}

func TestRoomDelete(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status":            200,
			"IsAdmin":           true,
			"DeleteRoomCalled":  1,
			"DeleteRoomSuccess": true,
		},
		"validation error (not admin)": {
			"status":            400,
			"IsAdmin":           false,
			"DeleteRoomCalled":  0,
			"DeleteRoomSuccess": false,
		},
		"failure to delete room": {
			"status":            500,
			"IsAdmin":           true,
			"DeleteRoomCalled":  1,
			"DeleteRoomSuccess": false,
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {

			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/room/:room_id", nil)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("uuid", "test-uuid-1234")
			c.Set("is_admin", expect["IsAdmin"].(bool))

			c.SetParamNames("room_id")
			c.SetParamValues("test-room-id")

			dto := dto.NewRoomDtoStruct()
			mongoSvcMock := new(mongo_svc_mock.RoomSvcMock)
			roomSvcMock := new(svc_mock.RoomSvcMock)

			var returnErr error = nil
			if !expect["DeleteRoomSuccess"].(bool) {
				returnErr = fmt.Errorf("DeleteRoom error")
			}
			mongoSvcMock.On("DeleteRoom", "test-room-id", mock.Anything).Return(returnErr).Times(expect["DeleteRoomCalled"].(int))

			handler := NewRoomHandler(mongoSvcMock, roomSvcMock, dto)
			err := handler.Delete(c)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["DeleteRoomCalled"].(int) != 0 {
				mongoSvcMock.AssertExpectations(t)
			} else {
				mongoSvcMock.AssertNotCalled(t, "DeleteRoom")
			}
		})
	}
}

func TestRoomAddMember(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status":          200,
			"IsAdmin":         true,
			"JoinRoomCalled":  1,
			"JoinRoomSuccess": true,
			"member_id":       "new-member-uuid-5678",
		},
		"validation error (not admin)": {
			"status":          400,
			"IsAdmin":         false,
			"JoinRoomCalled":  0,
			"JoinRoomSuccess": false,
			"member_id":       "new-member-uuid-5678",
		},
		"validation error (missing member_id)": {
			"status":          400,
			"IsAdmin":         true,
			"JoinRoomCalled":  0,
			"JoinRoomSuccess": false,
			"member_id":       "",
		},
		"failure to add member": {
			"status":          500,
			"IsAdmin":         true,
			"JoinRoomCalled":  1,
			"JoinRoomSuccess": false,
			"member_id":       "new-member-uuid-5678",
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {

			e := echo.New()
			e.Validator = &usecase.CustomValidator{Validator: validator.New()}

			body := map[string]interface{}{
				"member_id": expect["member_id"].(string),
			}

			jsonBody, _ := json.Marshal(body)
			reqBody := strings.NewReader(string(jsonBody))

			req := httptest.NewRequest(http.MethodPost, "/room/:room_id/add_member", reqBody)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("uuid", "test-uuid-1234")
			c.Set("is_admin", expect["IsAdmin"].(bool))

			c.SetParamNames("room_id")
			c.SetParamValues("test-room-id")

			dto := dto.NewRoomDtoStruct()
			mongoSvcMock := new(mongo_svc_mock.RoomSvcMock)
			roomSvcMock := new(svc_mock.RoomSvcMock)

			var returnErr error = nil
			if !expect["JoinRoomSuccess"].(bool) {
				returnErr = fmt.Errorf("JoinRoom error")
			}
			if expect["JoinRoomCalled"].(int) != 0 {
				mongoSvcMock.On("JoinRoom", "test-room-id", expect["member_id"].(string), mock.Anything).Return(returnErr).Times(expect["JoinRoomCalled"].(int))
			}

			handler := NewRoomHandler(mongoSvcMock, roomSvcMock, dto)
			err := handler.AddMember(c)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["JoinRoomCalled"].(int) != 0 {
				mongoSvcMock.AssertExpectations(t)
			} else {
				mongoSvcMock.AssertNotCalled(t, "JoinRoom")
			}
		})
	}
}

func TestRoomRemoveMember(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status":           200,
			"IsAdmin":          true,
			"LeaveRoomCalled":  1,
			"LeaveRoomSuccess": true,
			"member_id":        "member-uuid-5678",
		},
		"validation error (not admin)": {
			"status":           400,
			"IsAdmin":          false,
			"LeaveRoomCalled":  0,
			"LeaveRoomSuccess": false,
			"member_id":        "member-uuid-5678",
		},
		"validation error (missing member_id)": {
			"status":           400,
			"IsAdmin":          true,
			"LeaveRoomCalled":  0,
			"LeaveRoomSuccess": false,
			"member_id":        "",
		},
		"failure to remove member": {
			"status":           500,
			"IsAdmin":          true,
			"LeaveRoomCalled":  1,
			"LeaveRoomSuccess": false,
			"member_id":        "member-uuid-5678",
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {

			e := echo.New()
			e.Validator = &usecase.CustomValidator{Validator: validator.New()}

			body := map[string]interface{}{
				"member_id": expect["member_id"].(string),
			}

			jsonBody, _ := json.Marshal(body)
			reqBody := strings.NewReader(string(jsonBody))

			req := httptest.NewRequest(http.MethodPost, "/room/:room_id/remove_member", reqBody)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("uuid", "test-uuid-1234")
			c.Set("is_admin", expect["IsAdmin"].(bool))

			c.SetParamNames("room_id")
			c.SetParamValues("test-room-id")

			dto := dto.NewRoomDtoStruct()
			mongoSvcMock := new(mongo_svc_mock.RoomSvcMock)
			roomSvcMock := new(svc_mock.RoomSvcMock)

			var returnErr error = nil
			if !expect["LeaveRoomSuccess"].(bool) {
				returnErr = fmt.Errorf("LeaveRoom error")
			}
			if expect["LeaveRoomCalled"].(int) != 0 {
				mongoSvcMock.On("LeaveRoom", "test-room-id", expect["member_id"].(string), mock.Anything).Return(returnErr).Times(expect["LeaveRoomCalled"].(int))
			}

			handler := NewRoomHandler(mongoSvcMock, roomSvcMock, dto)
			err := handler.RemoveMember(c)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["LeaveRoomCalled"].(int) != 0 {
				mongoSvcMock.AssertExpectations(t)
			} else {
				mongoSvcMock.AssertNotCalled(t, "LeaveRoom")
			}
		})
	}
}
