package room_routing

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/handler"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/labstack/echo/v4"
)

type GroupRoutingInterface interface {
	DetailRoute() *echo.Group
	AdminRoute() *echo.Group
}

type GroupRouting struct {
	handler    handler.RoomHandlerInterface
	middleware *middleware.Middleware
	group      *echo.Group
	schema     string
}

func NewGroupRouting(
	handler handler.RoomHandlerInterface,
	middleware *middleware.Middleware,
	group *echo.Group,
	schema string,
) *GroupRouting {
	return &GroupRouting{
		handler:    handler,
		middleware: middleware,
		group:      group,
		schema:     schema,
	}
}
