package middleware

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestBeforeHandler(t *testing.T) {
	called := false
	beforeFunc := func(c echo.Context) error {
		called = true
		return nil
	}

	mw := BeforeHandler(beforeFunc)

	handler := mw(func(c echo.Context) error {
		return nil
	})

	err := handler(echo.New().NewContext(nil, nil))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !called {
		t.Errorf("Expected before function to be called")
	}
}

func TestBeforeHandlerForError(t *testing.T) {
	expectedErr := echo.NewHTTPError(400, "bad request")
	beforeFunc := func(c echo.Context) error {
		return expectedErr
	}

	mw := BeforeHandler(beforeFunc)

	handler := mw(func(c echo.Context) error {
		return nil
	})

	err := handler(echo.New().NewContext(nil, nil))
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestAfterHandler(t *testing.T) {
	called := false
	afterFunc := func(c echo.Context) error {
		called = true
		return nil
	}

	mw := AfterHandler(afterFunc)

	handler := mw(func(c echo.Context) error {
		return nil
	})

	err := handler(echo.New().NewContext(nil, nil))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !called {
		t.Errorf("Expected after function to be called")
	}
}

func TestAfterHandlerForError(t *testing.T) {
	expectedErr := echo.NewHTTPError(500, "internal server error")
	afterFunc := func(c echo.Context) error {
		return nil
	}

	mw := AfterHandler(afterFunc)

	handler := mw(func(c echo.Context) error {
		return expectedErr
	})

	err := handler(echo.New().NewContext(nil, nil))
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}
