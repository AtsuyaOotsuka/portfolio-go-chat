package handler_mock

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MockRoomHandler struct{}

func (h *MockRoomHandler) List(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"rooms": "list"})
}

func (h *MockRoomHandler) Create(c echo.Context) error {
	return c.JSON(http.StatusCreated, echo.Map{"room": "created"})
}

func (h *MockRoomHandler) Detail(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"room": "detail"})
}

func (h *MockRoomHandler) Members(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"members": "list"})
}

func (h *MockRoomHandler) Join(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"room": "joined"})
}

func (h *MockRoomHandler) Leave(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"room": "left"})
}

func (h *MockRoomHandler) Delete(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"room": "deleted"})
}

func (h *MockRoomHandler) AddMember(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"member": "added"})
}

func (h *MockRoomHandler) RemoveMember(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"member": "removed"})
}
