package room_routing

import "github.com/labstack/echo/v4"

func (r GroupRouting) AdminRoute() *echo.Group {
	roomAdminGroup := r.group.Group(
		r.schema,
		r.middleware.Room,
	)
	roomAdminGroup.DELETE("/delete", r.handler.Delete)
	roomAdminGroup.POST("/add_member", r.handler.AddMember)
	roomAdminGroup.DELETE("/remove_member", r.handler.RemoveMember)
	return roomAdminGroup
}
