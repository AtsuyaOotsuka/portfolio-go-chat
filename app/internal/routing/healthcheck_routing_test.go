package routing

import (
	"net/http"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/labstack/echo/v4"
)

type MockHealthCheckHandler struct{}

func (m *MockHealthCheckHandler) Check(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func TestHealthCheckRouting(t *testing.T) {
	expected := []funcs.ExpectedRoute{
		{Path: "/healthcheck", Method: "GET"},
	}

	e := echo.New()
	r := NewRouting(e, nil)
	r.HealthCheckRoute(&MockHealthCheckHandler{})

	funcs.EachExepectedRoute(expected, e, t)
}
