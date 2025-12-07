package routing

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/handler_mock"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/middleware_mock"
	"github.com/labstack/echo/v4"
)

func TestNewGroupRouting(t *testing.T) {
	e := echo.New()
	group := e.Group("/room")
	schema := "/:room_id"

	routing := NewGroupRouting(group, schema)

	if routing.group != group {
		t.Errorf("Expected group to be %v, got %v", group, routing.group)
	}
	if routing.schema != schema {
		t.Errorf("Expected schema to be %s, got %s", schema, routing.schema)
	}
}

func TestRoomRouting(t *testing.T) {
	expected := []funcs.ExpectedRoute{
		{Path: "/room/list", Method: "GET"},
		{Path: "/room/create", Method: "POST"},
		{Path: "/room/:room_id/detail", Method: "GET"},
		{Path: "/room/:room_id/members", Method: "GET"},
		{Path: "/room/:room_id/join", Method: "POST"},
		{Path: "/room/:room_id/leave", Method: "POST"},
		{Path: "/room/:room_id/admin/delete", Method: "DELETE"},
		{Path: "/room/:room_id/admin/add_member", Method: "POST"},
		{Path: "/room/:room_id/admin/remove_member", Method: "DELETE"},
	}

	e := echo.New()
	mw := &middleware.Middleware{
		RoomMV:    (&middleware_mock.MockRoomMiddleware{}).RoomMV,
		RoomAdmin: (&middleware_mock.MockRoomMiddleware{}).RoomAdmin,
	}
	r := NewRouting(e, mw)
	r.RoomRoute(&handler_mock.MockRoomHandler{})

	funcs.EachExepectedRoute(expected, e, t)
}
