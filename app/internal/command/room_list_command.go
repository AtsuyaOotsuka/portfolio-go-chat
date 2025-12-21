package command

import (
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/cmd_svc"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
)

type RoomListCommandInterface interface {
	SetUp(mongo usecase.MongoUseCaseInterface)
	Run()
}

type RoomListCommand struct {
	BaseCommand
	room_svc cmd_svc.RoomSvcInterface
}

func NewRoomListCommand() *RoomListCommand {
	return &RoomListCommand{}
}

func (c *RoomListCommand) SetUp(mongo usecase.MongoUseCaseInterface) {
	c.room_svc = cmd_svc.NewRoomSvcStruct(
		mongo,
	)
}

func (c *RoomListCommand) Run() {
	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	rooms, err := c.room_svc.ListRooms(ctx)
	if err != nil {
		fmt.Println("Error fetching rooms:", err.Error())
		return
	}

	for _, room := range rooms {
		fmt.Println("Room ID:", room.ID.Hex(), "Name:", room.Name)
	}
}
