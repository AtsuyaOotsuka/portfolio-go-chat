package middleware

import (
	"github.com/labstack/echo/v4"
)

type Middleware struct {
	Csrf echo.MiddlewareFunc
	Jwt  echo.MiddlewareFunc
	Room echo.MiddlewareFunc
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
