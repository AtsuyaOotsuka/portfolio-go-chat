package handler_mock

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MockMessageHandler struct{}

func (h *MockMessageHandler) List(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"messages": "list"})
}

func (h *MockMessageHandler) Send(c echo.Context) error {
	return c.JSON(http.StatusCreated, echo.Map{"message": "sent"})
}

func (h *MockMessageHandler) Read(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "read"})
}

func (h *MockMessageHandler) Delete(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "deleted"})
}
