package middleware

import "github.com/labstack/echo/v4"

type RoomAdminMiddlewareInterface interface {
	Handler() echo.MiddlewareFunc
}

type RoomAdminMiddleware struct{}

func NewRoomAdminMiddleware() RoomAdminMiddlewareInterface {
	return &RoomAdminMiddleware{}
}

func (m *RoomAdminMiddleware) Handler() echo.MiddlewareFunc {
	return BeforeHandler(func(c echo.Context) error {
		return nil
	})
}
