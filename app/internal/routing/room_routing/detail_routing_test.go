package room_routing

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/handler_mock"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/middleware_mock"
	"github.com/labstack/echo/v4"
)

func TestDetailRoute(t *testing.T) {
	expected := []funcs.ExpectedRoute{
		{Path: "/room/:room_id/detail", Method: "GET"},
		{Path: "/room/:room_id/members", Method: "GET"},
		{Path: "/room/:room_id/join", Method: "POST"},
		{Path: "/room/:room_id/leave", Method: "POST"},
	}
	e := echo.New()
	mw := &middleware.Middleware{
		RoomMV:    (&middleware_mock.MockRoomMiddleware{}).RoomMV,
		RoomAdmin: (&middleware_mock.MockRoomMiddleware{}).RoomAdmin,
	}

	r := NewGroupRouting(
		&handler_mock.MockRoomHandler{},
		mw,
		e.Group("/room"),
		"/:room_id",
	)
	r.DetailRoute()

	funcs.EachExepectedRoute(expected, e, t)
}
