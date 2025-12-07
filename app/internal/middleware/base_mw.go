package middleware

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
	"github.com/labstack/echo/v4"
)

type Middleware struct {
	e         *echo.Echo
	Csrf      echo.MiddlewareFunc
	Jwt       echo.MiddlewareFunc
	RoomMV    echo.MiddlewareFunc
	RoomAdmin echo.MiddlewareFunc
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
	roomMV := NewRoomMVMiddleware()
	roomAdmin := NewRoomAdminMiddleware()

	return &Middleware{
		e:         e,
		Csrf:      csrf.Handler(),
		Jwt:       jwt.Handler(),
		RoomMV:    roomMV.Handler(),
		RoomAdmin: roomAdmin.Handler(),
	}
}

func BeforeHandler(
	before func(c echo.Context) error,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			err := before(c)

			if err != nil {
				return err
			}

			return next(c)
		}
	}
}

func AfterHandler(
	after func(c echo.Context) error,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			err := next(c)

			if err != nil {
				return err
			}

			return after(c)
		}
	}
}
