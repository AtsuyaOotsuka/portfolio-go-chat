package middleware

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
	"github.com/labstack/echo/v4"
)

type Middleware struct {
	e    *echo.Echo
	Csrf echo.MiddlewareFunc
	Jwt  echo.MiddlewareFunc
}

func NewMiddleware(e *echo.Echo) *Middleware {

	csrf := NewCSRFMiddleware(
		service.NewCsrfSvcStruct(
			atylabcsrf.NewCsrfPkgStruct(),
		),
	)
	jwt := NewJWTMiddleware(
		atylabjwt.NewJwtSvc(),
	)

	return &Middleware{
		e:    e,
		Csrf: csrf.Handler(),
		Jwt:  jwt.Handler(),
	}
}
