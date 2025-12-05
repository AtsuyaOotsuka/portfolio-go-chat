package handler

import (
	"github.com/labstack/echo/v4"
)

type HealthCheckHandlerInterface interface {
	Check(c echo.Context) error
}

type HealthCheckHandler struct {
	BaseHandler
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Check(c echo.Context) error {
	return c.JSON(200, echo.Map{
		"status": "ok",
		"uuid":   c.Get("uuid"),
		"email":  c.Get("email"),
	})
}
