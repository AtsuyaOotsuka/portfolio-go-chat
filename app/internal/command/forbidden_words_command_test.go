package command

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/svc_mock/cmd_svc_mock"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var rooms = []model.Room{
	{
		ID:   primitive.NewObjectID(),
		Name: "General",
	},
	{
		ID:   primitive.NewObjectID(),
		Name: "Random",
	},
}

var messages = []model.Message{
	{
		ID:      primitive.NewObjectID(),
		RoomID:  rooms[0].ID.Hex(),
		Message: "This is a clean message.",
	},
	{
		ID:      primitive.NewObjectID(),
		RoomID:  rooms[0].ID.Hex(),
		Message: "This message contains badword2 a forbidden word.",
	},
	{
		ID:      primitive.NewObjectID(),
		RoomID:  rooms[1].ID.Hex(),
		Message: "Another clean message here.",
	},
	{
		ID:      primitive.NewObjectID(),
		RoomID:  rooms[1].ID.Hex(),
		Message: "Forbidden content badword1 in this one.",
	},
}

func TestForbiddenWordsCmdSetUp(t *testing.T) {
	cmd := NewForbiddenWordsCommand()
	cmd.SetUp(&usecase.MongoUseCaseStruct{}, 150)
	if cmd.room_svc == nil {
		t.Error("room_svc should not be nil after SetUp")
	}
	if cmd.timeOut != 150 {
		t.Errorf("timeOut should be 150 after SetUp, got %d", cmd.timeOut)
	}
}

func TestForbiddenWordsCmdRun(t *testing.T) {
	expected := map[string]map[string]any{
		"success": {
			"ListRoomsError": nil,
		},
		"error": {
			"ListRoomsError": errors.New("failed to list rooms"),
		},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {

			roomSvcMock := new(cmd_svc_mock.RoomSvcMock)
			roomSvcMock.On("ListRooms", mock.Anything).Return(rooms, expect["ListRoomsError"])
			messageSvcMock := new(cmd_svc_mock.MessageSvcMock)
			if expect["ListRoomsError"] == nil {
				messageSvcMock.On("GetMessageList", rooms[0].ID.Hex(), mock.Anything).Return([]model.Message{messages[0], messages[1]}, nil)
				messageSvcMock.On("GetMessageList", rooms[1].ID.Hex(), mock.Anything).Return([]model.Message{messages[2], messages[3]}, nil)
				messageSvcMock.On("ContainsForbiddenWords", messages[0].Message).Return(false)
				messageSvcMock.On("ContainsForbiddenWords", messages[1].Message).Return(true)
				messageSvcMock.On("ContainsForbiddenWords", messages[2].Message).Return(false)
				messageSvcMock.On("ContainsForbiddenWords", messages[3].Message).Return(true)
			}

			cmd := NewForbiddenWordsCommand()
			cmd.room_svc = roomSvcMock
			cmd.message_svc = messageSvcMock
			cmd.timeOut = 100

			outPut := funcs.CaptureStdout(t, func() {
				cmd.Run([]string{})
			})

			roomSvcMock.AssertExpectations(t)
			messageSvcMock.AssertExpectations(t)

			if expect["ListRoomsError"] != nil {
				return
			}

			targetIds := []string{
				messages[1].ID.Hex(),
				messages[3].ID.Hex(),
			}
			unTargetIds := []string{
				messages[0].ID.Hex(),
				messages[2].ID.Hex(),
			}

			for i, targetId := range targetIds {
				if !strings.Contains(outPut, targetId) {
					t.Errorf("Expected output to contain targetId %s, but it did not. Index: %d", targetId, i)
				}
			}

			for i, unTargetId := range unTargetIds {
				if strings.Contains(outPut, unTargetId) {
					t.Errorf("Expected output to NOT contain unTargetId %s, but it did. Index: %d", unTargetId, i)
				}
			}

			if strings.Count(outPut, "Forbidden word found in Room ID:") != 2 {
				t.Errorf("Expected output to contain 2 forbidden word findings, but got %d", strings.Count(outPut, "Forbidden word found in Room ID:"))
			}

			if !strings.Contains(outPut, "処理完了") {
				t.Error("Expected output to contain '処理完了', but it did not.")
			}

		})
	}
}

func TestForbiddenWordsCmdRunGetMessageListError(t *testing.T) {
	roomSvcMock := new(cmd_svc_mock.RoomSvcMock)
	roomSvcMock.On("ListRooms", mock.Anything).Return([]model.Room{
		rooms[0],
	}, nil)
	messageSvcMock := new(cmd_svc_mock.MessageSvcMock)
	messageSvcMock.On("GetMessageList", rooms[0].ID.Hex(), mock.Anything).Return([]model.Message{}, errors.New("failed to get message list"))

	cmd := NewForbiddenWordsCommand()
	cmd.room_svc = roomSvcMock
	cmd.message_svc = messageSvcMock
	cmd.timeOut = 100

	outPut := funcs.CaptureStdout(t, func() {
		cmd.Run([]string{})
	})

	if !strings.Contains(outPut, "Error processing messages") {
		t.Error("Expected error message for listing rooms, but it was not found in output.")
	}

	roomSvcMock.AssertExpectations(t)
	messageSvcMock.AssertExpectations(t)
}

func TestForbiddenWordsCmdRunGroupError(t *testing.T) {
	roomSvcMock := new(cmd_svc_mock.RoomSvcMock)
	roomSvcMock.On("ListRooms", mock.Anything).Return(rooms, nil)
	messageSvcMock := new(cmd_svc_mock.MessageSvcMock)
	messageSvcMock.On("GetMessageList", rooms[0].ID.Hex(), mock.Anything).
		Return([]model.Message{messages[0], messages[1]}, nil)

	messageSvcMock.On("GetMessageList", rooms[1].ID.Hex(), mock.Anything).
		Return([]model.Message{messages[2], messages[3]}, nil)

	cmd := NewForbiddenWordsCommand()
	cmd.room_svc = roomSvcMock
	cmd.message_svc = messageSvcMock
	cmd.timeOut = 0

	outPut := funcs.CaptureStdout(t, func() {
		cmd.Run([]string{})
	})

	if !strings.Contains(outPut, "Error processing messages:") {
		t.Error("Expected no error output, but got an error.")
	}

	roomSvcMock.AssertExpectations(t)
	messageSvcMock.AssertExpectations(t)
}

func TestForbiddenWordsCmdRunMultichan(t *testing.T) {

	roomSvcMock := new(cmd_svc_mock.RoomSvcMock)
	roomSvcMock.On("ListRooms", mock.Anything).Return(rooms, nil)
	messageSvcMock := new(cmd_svc_mock.MessageSvcMock)
	messageSvcMock.On("GetMessageList", mock.Anything, mock.Anything).
		Return([]model.Message{}, nil).
		Run(func(args mock.Arguments) {
			roomId := args.String(0)
			fmt.Println("start", roomId)
			time.Sleep(1 * time.Second)
			fmt.Println("end", roomId)
		}).
		Twice()

	cmd := NewForbiddenWordsCommand()
	cmd.room_svc = roomSvcMock
	cmd.message_svc = messageSvcMock
	cmd.timeOut = 100

	outPut := funcs.CaptureStdout(t, func() {
		cmd.Run([]string{})
	})
	lines := strings.Split(strings.TrimSpace(outPut), "\n")

	for i, line := range lines {
		switch i {
		case 0, 1:
			if !strings.Contains(line, "start") {
				t.Errorf("Expected 'start' in line %d, but got: %s", i, line)
			}
		case 2, 3:
			if !strings.Contains(line, "end") {
				t.Errorf("Expected 'end' in line %d, but got: %s", i, line)
			}
		}
	}

	roomSvcMock.AssertExpectations(t)
	messageSvcMock.AssertExpectations(t)
}
