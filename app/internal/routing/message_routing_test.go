package routing

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/handler_mock"
	"github.com/labstack/echo/v4"
)

func TestMessageRoute(t *testing.T) {
	expected := []funcs.ExpectedRoute{
		{Path: "/message/:room_id/list", Method: "GET"},
		{Path: "/message/:room_id/send", Method: "POST"},
		{Path: "/message/:room_id/read", Method: "POST"},
		{Path: "/message/:room_id/delete", Method: "DELETE"},
	}
	e := echo.New()
	mw := &middleware.Middleware{}
	r := NewRouting(e, mw)
	r.MessageRoute(&handler_mock.MockMessageHandler{})

	funcs.EachExepectedRoute(expected, e, t)

}
