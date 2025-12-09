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
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMessageList(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status":                200,
			"IsJoinedRoomCalled":    1,
			"IsJoinedRoomSuccess":   true,
			"GetMessageListCalled":  1,
			"GetMessageListSuccess": true,
			"success":               true,
		},
		"failure to check joined room": {
			"status":                500,
			"IsJoinedRoomCalled":    1,
			"IsJoinedRoomSuccess":   false,
			"GetMessageListCalled":  0,
			"GetMessageListSuccess": true,
			"success":               false,
		},
		"failure to get message list": {
			"status":                500,
			"IsJoinedRoomCalled":    1,
			"IsJoinedRoomSuccess":   true,
			"GetMessageListCalled":  1,
			"GetMessageListSuccess": false,
			"success":               false,
		},
	}

	var err error
	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/message/:room_id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("room_id")
			c.SetParamValues("test-room-id")
			c.Set("uuid", "test-uuid-1234")

			dto := dto.NewMessageDtoStruct()
			roomSvcMock := new(svc_mock.RoomSvcMock)
			messageSvcMock := new(svc_mock.MessageSvcMock)

			returnData := []model.Message{
				{
					ID:            primitive.NewObjectID(),
					RoomID:        "test-room-id",
					Sender:        "test-sender-id",
					Message:       "test message",
					CreatedAt:     time.Now(),
					IsReadUserIds: []string{"test-sender-id", "test-uuid-1234"},
				},
				{
					ID:            primitive.NewObjectID(),
					RoomID:        "test-room-id",
					Sender:        "another-sender-id",
					Message:       "another test message",
					CreatedAt:     time.Now(),
					IsReadUserIds: []string{"another-sender-id"},
				},
			}

			var isJoinedErr error = nil
			if !expect["IsJoinedRoomSuccess"].(bool) {
				isJoinedErr = assert.AnError
			}
			var getMessageListErr error = nil
			if !expect["GetMessageListSuccess"].(bool) {
				getMessageListErr = assert.AnError
				returnData = []model.Message{}
			}

			if expect["IsJoinedRoomCalled"].(int) > 0 {
				roomSvcMock.
					On("IsJoinedRoom", "test-room-id", "test-uuid-1234").
					Return(isJoinedErr).
					Times(expect["IsJoinedRoomCalled"].(int))
			}

			if expect["GetMessageListCalled"].(int) > 0 {
				messageSvcMock.
					On("GetMessageList", "test-room-id").
					Return(returnData, getMessageListErr).
					Times(expect["GetMessageListCalled"].(int))
			}

			handler := NewMessageHandler(messageSvcMock, roomSvcMock, dto)
			err = handler.List(c)

			assert.NoError(t, err)

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["IsJoinedRoomCalled"].(int) > 0 {
				roomSvcMock.AssertExpectations(t)
			} else {
				roomSvcMock.AssertNotCalled(t, "IsJoinedRoom")
			}

			if expect["GetMessageListCalled"].(int) > 0 {
				messageSvcMock.AssertExpectations(t)
			} else {
				messageSvcMock.AssertNotCalled(t, "GetMessageList")
			}

			if expect["status"].(int) != http.StatusOK {
				return
			}

			result := map[string][]interface{}{}
			err = json.Unmarshal(rec.Body.Bytes(), &result)
			fmt.Println(rec.Body.String())
			assert.NoError(t, err)

			assert.Len(t, result["messages"], 2)

			assert.Equal(t, "test message", result["messages"][0].(map[string]interface{})["Message"])
			assert.Equal(t, "test-sender-id", result["messages"][0].(map[string]interface{})["Sender"])
			assert.True(t, result["messages"][0].(map[string]interface{})["IsRead"].(bool))

			assert.Equal(t, "another test message", result["messages"][1].(map[string]interface{})["Message"])
			assert.Equal(t, "another-sender-id", result["messages"][1].(map[string]interface{})["Sender"])
			assert.False(t, result["messages"][1].(map[string]interface{})["IsRead"].(bool))
		})
	}
}

func TestMessageSend(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status": 200,
			"body": map[string]interface{}{
				"message": "Hello, world!",
			},
			"success":             true,
			"IsJoinedRoomCalled":  1,
			"IsJoinedRoomSuccess": true,
			"SendMessageCalled":   1,
			"SendMessageSuccess":  true,
		},
		"validation error (missing message)": {
			"status":              400,
			"body":                map[string]interface{}{},
			"success":             false,
			"IsJoinedRoomCalled":  0,
			"IsJoinedRoomSuccess": true,
			"SendMessageCalled":   0,
			"SendMessageSuccess":  true,
		},
		"failure to check joined room": {
			"status": 500,
			"body": map[string]interface{}{
				"message": "Hello, world!",
			},
			"success":             false,
			"IsJoinedRoomCalled":  1,
			"IsJoinedRoomSuccess": false,
			"SendMessageCalled":   0,
			"SendMessageSuccess":  true,
		},
		"failure to send message": {
			"status": 500,
			"body": map[string]interface{}{
				"message": "Hello, world!",
			},
			"success":             false,
			"IsJoinedRoomCalled":  1,
			"IsJoinedRoomSuccess": true,
			"SendMessageCalled":   1,
			"SendMessageSuccess":  false,
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &usecase.CustomValidator{Validator: validator.New()}

			body := expect["body"].(map[string]interface{})
			jsonBody, _ := json.Marshal(body)
			reqBody := strings.NewReader(string(jsonBody))

			req := httptest.NewRequest(http.MethodPost, "/message/:room_id/send", reqBody)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("room_id")
			c.SetParamValues("test-room-id")
			c.Set("uuid", "test-uuid-1234")

			dto := dto.NewMessageDtoStruct()
			roomSvcMock := new(svc_mock.RoomSvcMock)
			messageSvcMock := new(svc_mock.MessageSvcMock)

			var isJoinedErr error = nil
			if !expect["IsJoinedRoomSuccess"].(bool) {
				isJoinedErr = assert.AnError
			}
			var sendMessageErr error = nil
			if !expect["SendMessageSuccess"].(bool) {
				sendMessageErr = assert.AnError
			}

			if expect["IsJoinedRoomCalled"].(int) > 0 {
				roomSvcMock.
					On("IsJoinedRoom", "test-room-id", "test-uuid-1234").
					Return(isJoinedErr).
					Times(expect["IsJoinedRoomCalled"].(int))
			}

			if expect["SendMessageCalled"].(int) > 0 {
				messageSvcMock.
					On("SendMessage", mock.AnythingOfType("model.Message")).
					Return("new-message-id-5678", sendMessageErr).
					Times(expect["SendMessageCalled"].(int))
			}

			handler := NewMessageHandler(messageSvcMock, roomSvcMock, dto)
			err := handler.Send(c)

			assert.NoError(t, err)

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["IsJoinedRoomCalled"].(int) > 0 {
				roomSvcMock.AssertExpectations(t)
			} else {
				roomSvcMock.AssertNotCalled(t, "IsJoinedRoom")
			}

			if expect["SendMessageCalled"].(int) > 0 {
				messageSvcMock.AssertExpectations(t)
			} else {
				messageSvcMock.AssertNotCalled(t, "SendMessage")
			}

			if expect["status"].(int) != http.StatusOK {
				return
			}

			result := map[string]interface{}{}
			err = json.Unmarshal(rec.Body.Bytes(), &result)
			assert.NoError(t, err)

			assert.Equal(t, "new-message-id-5678", result["message_id"])
		})
	}
}
