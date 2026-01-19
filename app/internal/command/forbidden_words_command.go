package command

import (
	"context"
	"fmt"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/cmd_svc"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
	"golang.org/x/sync/errgroup"
)

type ForbiddenWordsCommandInterface interface {
	SetUp(mongo usecase.MongoUseCaseInterface, timeOut int)
	Run(args []string)
}

type ForbiddenWordsCommand struct {
	BaseCommand
	room_svc    cmd_svc.RoomSvcInterface
	message_svc cmd_svc.MessageSvcInterface
	timeOut     int
}

func NewForbiddenWordsCommand() *ForbiddenWordsCommand {
	return &ForbiddenWordsCommand{}
}

func (c *ForbiddenWordsCommand) SetUp(
	mongo usecase.MongoUseCaseInterface,
	timeOut int,
) {
	c.room_svc = cmd_svc.NewRoomSvcStruct(
		mongo,
	)
	c.message_svc = cmd_svc.NewMessageSvcStruct(
		mongo,
	)
	c.timeOut = timeOut
}

func (c *ForbiddenWordsCommand) Run(args []string) {
	// 全体のタイムアウトを100秒に設定
	gctx, gctxCancel := context.WithTimeout(context.Background(), time.Duration(c.timeOut)*time.Second)
	defer gctxCancel()

	g, gctx := errgroup.WithContext(gctx)

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	rooms, err := c.room_svc.ListRooms(ctx)
	if err != nil {
		fmt.Println("Error fetching rooms:", err.Error())
		return
	}

	for _, room := range rooms {
		g.Go(func() error {
			room := room // クロージャ内で正しいroomを参照するために変数を再定義
			roomId := room.ID.Hex()

			mctx := atylabmongo.NewMongoCtxSvc()
			defer mctx.Cancel()
			messageList, err := c.message_svc.GetMessageList(roomId, mctx)

			if err != nil {
				return err
			}
			for _, message := range messageList {
				if err := gctx.Err(); err != nil {
					return err
				}
				if c.message_svc.ContainsForbiddenWords(message.Message) {
					fmt.Printf("Forbidden word found in Room ID: %s, Message ID: %s, Content: %s\n", roomId, message.ID.Hex(), message.Message)
				}
			}
			// すべてのメッセージを処理した後の処理
			return nil
		})
	}

	// goroutineの完了を待つ
	if err := g.Wait(); err != nil {
		fmt.Println("Error processing messages:", err.Error())
		gctxCancel()
		return
	}

	fmt.Println("処理完了")
}
