package middleware

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestNewMiddleware(t *testing.T) {
	e := echo.New()
	mv := NewMiddleware(e)

	if mv.e != e {
		t.Errorf("Expected echo instance to be set correctly")
	}
	if mv.Csrf == nil {
		t.Errorf("Expected Csrf middleware to be initialized")
	}
	if mv.Jwt == nil {
		t.Errorf("Expected Jwt middleware to be initialized")
	}
}
