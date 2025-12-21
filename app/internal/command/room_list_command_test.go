package command

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/svc_mock/cmd_svc_mock"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRoomListCmdSetUp(t *testing.T) {
	cmd := NewRoomListCommand()
	cmd.SetUp(&usecase.MongoUseCaseStruct{})
	if cmd.room_svc == nil {
		t.Error("room_svc should not be nil after SetUp")
	}
}

func TestRoomListCmdRun(t *testing.T) {
	rooms := []model.Room{
		{
			ID:   primitive.NewObjectID(),
			Name: "General",
		},
		{
			ID:   primitive.NewObjectID(),
			Name: "Random",
		},
	}

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

			// モックのMongoUseCaseInterfaceを作成
			svcMock := new(cmd_svc_mock.RoomSvcMock)
			svcMock.On("ListRooms", mock.Anything).Return(rooms, expect["ListRoomsError"])

			// RoomListCommandのインスタンスを作成
			cmd := NewRoomListCommand()
			cmd.room_svc = svcMock

			// 標準出力をキャプチャ
			outPut := funcs.CaptureStdout(t, func() {
				cmd.Run()
			})

			// 期待される出力を定義
			if expect["ListRoomsError"] != nil {
				expectedError := "Error fetching rooms: failed to list rooms"
				if strings.TrimSpace(outPut) != expectedError {
					t.Errorf("expected %q but got %q", expectedError, strings.TrimSpace(outPut))
				}
				return
			}

			result := strings.Split(outPut, "\n")
			expectedOutputs := []string{
				fmt.Sprintf("Room ID: %s Name: General", rooms[0].ID.Hex()),
				fmt.Sprintf("Room ID: %s Name: Random", rooms[1].ID.Hex()),
				"",
			}

			// 出力を検証
			for i, expected := range expectedOutputs {
				if result[i] != expected {
					t.Errorf("expected %q but got %q", expected, result[i])
				}
			}

		})
	}
}
