package middleware

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/consts"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"github.com/labstack/echo/v4"
)

type RoomMVMiddlewareInterface interface {
	Handler() echo.MiddlewareFunc
}

type RoomMVMiddleware struct {
	roomSvc service.RoomSvcInterface
}

func NewRoomMiddleware(
	roomSvc service.RoomSvcInterface,
) RoomMVMiddlewareInterface {
	return &RoomMVMiddleware{
		roomSvc: roomSvc,
	}
}

func (m *RoomMVMiddleware) Handler() echo.MiddlewareFunc {
	return BeforeHandler(func(c echo.Context) error {
		ctx := atylabmongo.NewMongoCtxSvc()
		defer ctx.Cancel()

		roomID := c.Param("room_id")
		uuid := c.Get(consts.ContextKeys.Uuid).(string)
		room, err := m.roomSvc.GetRoom(roomID, ctx)
		if err != nil {
			return echo.NewHTTPError(404, "room not found")
		}
		c.Set(consts.ContextKeys.RoomModel, room)
		c.Set(consts.ContextKeys.IsAdmin, m.roomSvc.IsOwner(room, uuid))
		c.Set(consts.ContextKeys.IsMember, m.roomSvc.IsMember(room, uuid))

		return nil
	})
}
