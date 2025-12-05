package funcs

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type ExpectedRoute struct {
	Path   string
	Method string
}

func EachExepectedRoute(
	expected []ExpectedRoute,
	e *echo.Echo,
	t assert.TestingT,
) {
	for _, er := range expected {
		found := false
		for _, route := range e.Routes() {
			if route.Path == er.Path && route.Method == er.Method {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected route %s [%s] to be registered", er.Path, er.Method)
	}
}
