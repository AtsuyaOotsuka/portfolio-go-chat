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

func TestMessageList(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status":                200,
			"IsMember":              true,
			"GetMessageListCalled":  1,
			"GetMessageListSuccess": true,
			"success":               true,
		},
		"forbidden (not a member)": {
			"status":                403,
			"IsMember":              false,
			"GetMessageListCalled":  0,
			"GetMessageListSuccess": false,
			"success":               false,
		},
		"failure to get message list": {
			"status":                500,
			"IsMember":              true,
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
			c.Set("is_member", expect["IsMember"].(bool))

			dto := dto.NewMessageDtoStruct()
			messageSvcMock := new(mongo_svc_mock.MessageSvcMock)

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
			var getMessageListErr error = nil
			if !expect["GetMessageListSuccess"].(bool) {
				getMessageListErr = assert.AnError
				returnData = []model.Message{}
			}

			if expect["GetMessageListCalled"].(int) > 0 {
				messageSvcMock.
					On("GetMessageList", "test-room-id", mock.Anything).
					Return(returnData, getMessageListErr).
					Times(expect["GetMessageListCalled"].(int))
			}

			handler := NewMessageHandler(messageSvcMock, dto)
			err = handler.List(c)

			assert.NoError(t, err)

			assert.Equal(t, expect["status"].(int), rec.Code)

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
			"IsMember":           true,
			"success":            true,
			"SendMessageCalled":  1,
			"SendMessageSuccess": true,
		},
		"validation error (missing message)": {
			"status":             400,
			"body":               map[string]interface{}{},
			"IsMember":           false,
			"success":            false,
			"SendMessageCalled":  0,
			"SendMessageSuccess": true,
		},
		"forbidden (not a member)": {
			"status": 403,
			"body": map[string]interface{}{
				"message": "Hello, world!",
			},
			"IsMember":           false,
			"success":            false,
			"SendMessageCalled":  0,
			"SendMessageSuccess": true,
		},
		"failure to send message": {
			"status": 500,
			"body": map[string]interface{}{
				"message": "Hello, world!",
			},
			"IsMember":           true,
			"success":            false,
			"SendMessageCalled":  1,
			"SendMessageSuccess": false,
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
			c.Set("is_member", expect["IsMember"].(bool))

			dto := dto.NewMessageDtoStruct()

			messageSvcMock := new(mongo_svc_mock.MessageSvcMock)

			var sendMessageErr error = nil
			if !expect["SendMessageSuccess"].(bool) {
				sendMessageErr = assert.AnError
			}

			if expect["SendMessageCalled"].(int) > 0 {
				messageSvcMock.
					On("SendMessage", mock.AnythingOfType("model.Message"), mock.Anything).
					Return("new-message-id-5678", sendMessageErr).
					Times(expect["SendMessageCalled"].(int))
			}

			handler := NewMessageHandler(messageSvcMock, dto)
			err := handler.Send(c)

			assert.NoError(t, err)

			assert.Equal(t, expect["status"].(int), rec.Code)

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

func TestMessageRead(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status": 200,
			"body": map[string]interface{}{
				"message_ids": []string{"msgid1", "msgid2"},
			},
			"IsMember":            true,
			"ReadMessagesCalled":  1,
			"ReadMessagesSuccess": true,
		},
		"validation error (missing message_ids)": {
			"status":              400,
			"body":                map[string]interface{}{},
			"IsMember":            false,
			"ReadMessagesCalled":  0,
			"ReadMessagesSuccess": true,
		},
		"forbidden (not a member)": {
			"status": 403,
			"body": map[string]interface{}{
				"message_ids": []string{"msgid1", "msgid2"},
			},
			"IsMember":            false,
			"ReadMessagesCalled":  0,
			"ReadMessagesSuccess": true,
		},
		"failure to read messages": {
			"status": 500,
			"body": map[string]interface{}{
				"message_ids": []string{"msgid1", "msgid2"},
			},
			"IsMember":            true,
			"ReadMessagesCalled":  1,
			"ReadMessagesSuccess": false,
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &usecase.CustomValidator{Validator: validator.New()}

			body := expect["body"].(map[string]interface{})
			jsonBody, _ := json.Marshal(body)
			reqBody := strings.NewReader(string(jsonBody))
			req := httptest.NewRequest(http.MethodPost, "/message/:room_id/read", reqBody)

			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("room_id")
			c.SetParamValues("test-room-id")
			c.Set("uuid", "test-uuid-1234")
			c.Set("is_member", expect["IsMember"].(bool))

			dto := dto.NewMessageDtoStruct()

			messageSvcMock := new(mongo_svc_mock.MessageSvcMock)

			var readMessagesErr error = nil
			if !expect["ReadMessagesSuccess"].(bool) {
				readMessagesErr = assert.AnError
			}

			if expect["ReadMessagesCalled"].(int) > 0 {
				messageSvcMock.
					On("ReadMessages", []string{"msgid1", "msgid2"}, "test-room-id", "test-uuid-1234", mock.Anything).
					Return(readMessagesErr).
					Times(expect["ReadMessagesCalled"].(int))
			}

			handler := NewMessageHandler(messageSvcMock, dto)
			err := handler.Read(c)

			assert.NoError(t, err)

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["ReadMessagesCalled"].(int) > 0 {
				messageSvcMock.AssertExpectations(t)
			} else {
				messageSvcMock.AssertNotCalled(t, "ReadMessages")
			}

			if expect["status"].(int) != http.StatusOK {
				return
			}

			result := map[string]interface{}{}
			err = json.Unmarshal(rec.Body.Bytes(), &result)
			assert.NoError(t, err)

			assert.Equal(t, "success", result["status"])
		})
	}
}

func TestMessageDelete(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"status": 200,
			"body": map[string]interface{}{
				"message_id": "msgid1",
			},
			"IsMember":             true,
			"IsSenderCalled":       1,
			"IsSenderSuccess":      true,
			"IsOwner":              true,
			"DeleteMessageCalled":  1,
			"DeleteMessageSuccess": true,
		},
		"success by room owner": {
			"status": 200,
			"body": map[string]interface{}{
				"message_id": "msgid1",
			},
			"IsMember":             true,
			"IsSenderCalled":       1,
			"IsSenderSuccess":      true,
			"IsOwner":              true,
			"DeleteMessageCalled":  1,
			"DeleteMessageSuccess": true,
		},
		"failure to check is sender or room owner": {
			"status": 403,
			"body": map[string]interface{}{
				"message_id": "msgid1",
			},
			"IsMember":             true,
			"IsSenderCalled":       1,
			"IsSenderSuccess":      false,
			"IsOwner":              false,
			"DeleteMessageCalled":  0,
			"DeleteMessageSuccess": true,
		},
		"failure to delete message": {
			"status": 500,
			"body": map[string]interface{}{
				"message_id": "msgid1",
			},
			"IsMember":             true,
			"IsSenderCalled":       1,
			"IsSenderSuccess":      true,
			"IsOwner":              true,
			"DeleteMessageCalled":  1,
			"DeleteMessageSuccess": false,
		},
		"validation error (missing message_id)": {
			"status":               400,
			"body":                 map[string]interface{}{},
			"IsMember":             false,
			"IsSenderCalled":       0,
			"IsSenderSuccess":      true,
			"IsOwner":              true,
			"DeleteMessageCalled":  0,
			"DeleteMessageSuccess": true,
		},
		"forbidden (not a member)": {
			"status":               403,
			"body":                 map[string]interface{}{"message_id": "msgid1"},
			"IsMember":             false,
			"IsSenderCalled":       0,
			"IsSenderSuccess":      true,
			"IsOwner":              true,
			"DeleteMessageCalled":  0,
			"DeleteMessageSuccess": true,
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {
			e := echo.New()
			e.Validator = &usecase.CustomValidator{Validator: validator.New()}

			body := expect["body"].(map[string]interface{})
			jsonBody, _ := json.Marshal(body)
			reqBody := strings.NewReader(string(jsonBody))
			req := httptest.NewRequest(http.MethodPost, "/message/:room_id/read", reqBody)

			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.SetParamNames("room_id")
			c.SetParamValues("test-room-id")
			c.Set("uuid", "test-uuid-1234")
			c.Set("is_member", expect["IsMember"].(bool))
			c.Set("is_admin", expect["IsOwner"].(bool))

			dto := dto.NewMessageDtoStruct()

			messageSvcMock := new(mongo_svc_mock.MessageSvcMock)

			var isSenderErr error = nil
			if !expect["IsSenderSuccess"].(bool) {
				isSenderErr = assert.AnError
			}

			var deleteMessageErr error = nil
			if !expect["DeleteMessageSuccess"].(bool) {
				deleteMessageErr = assert.AnError
			}

			if expect["IsSenderCalled"].(int) > 0 {
				messageSvcMock.
					On("IsSender", "msgid1", "test-room-id", "test-uuid-1234", mock.Anything).
					Return(isSenderErr).
					Times(expect["IsSenderCalled"].(int))
			}

			if expect["DeleteMessageCalled"].(int) > 0 {
				messageSvcMock.
					On("DeleteMessage", "msgid1", "test-room-id", mock.Anything).
					Return(deleteMessageErr).
					Times(expect["DeleteMessageCalled"].(int))
			}

			handler := NewMessageHandler(messageSvcMock, dto)
			err := handler.Delete(c)

			assert.NoError(t, err)

			assert.Equal(t, expect["status"].(int), rec.Code)

			if expect["IsSenderCalled"].(int) > 0 {
				messageSvcMock.AssertExpectations(t)
			} else {
				messageSvcMock.AssertNotCalled(t, "IsSender")
			}

			if expect["DeleteMessageCalled"].(int) > 0 {
				messageSvcMock.AssertExpectations(t)
			} else {
				messageSvcMock.AssertNotCalled(t, "DeleteMessage")
			}

			if expect["status"].(int) != http.StatusOK {
				return
			}

			result := map[string]interface{}{}
			err = json.Unmarshal(rec.Body.Bytes(), &result)
			assert.NoError(t, err)

			assert.Equal(t, "success", result["status"])

		})
	}
}
