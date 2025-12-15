package room_routing

import (
	"github.com/labstack/echo/v4"
)

func (r GroupRouting) DetailRoute() *echo.Group {
	roomDetailGroup := r.group.Group(r.schema, r.middleware.Room)
	roomDetailGroup.POST("/join", r.handler.Join)
	roomDetailGroup.GET("/members", r.handler.Members)
	roomDetailGroup.POST("/leave", r.handler.Leave)

	return roomDetailGroup
}
