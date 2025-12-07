package middleware_mock

import "github.com/labstack/echo/v4"

type MockRoomMiddleware struct{}

func (m *MockRoomMiddleware) RoomMV(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func (m *MockRoomMiddleware) RoomAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
