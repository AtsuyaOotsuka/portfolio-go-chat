package middleware

import "github.com/labstack/echo/v4"

type RoomMVMiddlewareInterface interface {
	Handler() echo.MiddlewareFunc
}

type RoomMVMiddleware struct{}

func NewRoomMVMiddleware() RoomMVMiddlewareInterface {
	return &RoomMVMiddleware{}
}

func (m *RoomMVMiddleware) Handler() echo.MiddlewareFunc {
	return BeforeHandler(func(c echo.Context) error {
		return nil
	})
}
