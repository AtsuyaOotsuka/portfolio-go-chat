package middleware

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service"
	"github.com/labstack/echo/v4"
)

type CSRFMiddlewareInterface interface {
	Handler() echo.MiddlewareFunc
}

type CSRFMiddleware struct {
	csrf service.CsrfSvcInterface
}

func NewCSRFMiddleware(
	v service.CsrfSvcInterface,
) CSRFMiddlewareInterface {
	return &CSRFMiddleware{
		csrf: v,
	}
}

func (m *CSRFMiddleware) Handler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == echo.GET || c.Request().Method == echo.HEAD {
				return next(c)
			}
			token := c.Request().Header.Get("X-CSRF-Token")
			if token == "" {
				token = c.FormValue("_csrf")
			}
			if token == "" {
				cookie, err := c.Cookie("csrf_token")
				if err == nil {
					token, _ = url.QueryUnescape(cookie.Value)
				}
			}
			if token == "" {
				return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
					"error": "not set csrf token",
				})
			}
			if err := m.csrf.Verify(
				token,
				os.Getenv("CSRF_TOKEN"),
				time.Now().Unix(),
			); err != nil {
				return echo.NewHTTPError(http.StatusForbidden, echo.Map{
					"error": "invalid csrf token",
				})
			}
			return next(c)
		}
	}
}
