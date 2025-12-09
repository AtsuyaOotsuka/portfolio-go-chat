package routing

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/handler"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/routing/room_routing"
)

func (r *Routing) RoomRoute(
	handler handler.RoomHandlerInterface,
) {
	roomGroup := r.echo.Group("/room")

	roomGroup.GET("/list", handler.List)
	roomGroup.POST("/create", handler.Create)

	roomDetailGroup := room_routing.NewGroupRouting(
		handler,
		r.middleware,
		roomGroup,
		"/:room_id",
	).DetailRoute()

	roomAdminGroup := room_routing.NewGroupRouting(
		handler,
		r.middleware,
		roomDetailGroup,
		"/admin",
	).AdminRoute()

	r.Finalize(roomAdminGroup)
}
