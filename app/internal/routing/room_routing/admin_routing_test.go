package room_routing

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/handler_mock"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/middleware_mock"
	"github.com/labstack/echo/v4"
)

func TestAdminRoute(t *testing.T) {
	expected := []funcs.ExpectedRoute{
		{Path: "/room/:room_id/admin/delete", Method: "DELETE"},
		{Path: "/room/:room_id/admin/add_member", Method: "POST"},
		{Path: "/room/:room_id/admin/remove_member", Method: "DELETE"},
	}
	e := echo.New()
	mw := &middleware.Middleware{
		Room: (&middleware_mock.MockRoomMiddleware{}).RoomMV,
	}

	room_g := e.Group("/room")
	detail_g := room_g.Group("/:room_id")

	r := NewGroupRouting(
		&handler_mock.MockRoomHandler{},
		mw,
		detail_g,
		"/admin",
	)
	r.AdminRoute()

	funcs.EachExepectedRoute(expected, e, t)
}
