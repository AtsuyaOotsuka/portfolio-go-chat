package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/svc_mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTestCsrfHandler(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", "valid_token", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

	e := echo.New()
	e.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())

	e.POST("/test", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "not set csrf token")
}

func TestCSRFMiddleware_InvalidToken(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", "invalid-token", mock.Anything, mock.AnythingOfType("int64")).Return(fmt.Errorf("invalid"))

	e := echo.New()
	e.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())
	e.POST("/test", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("X-CSRF-Token", "invalid-token")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "invalid csrf token")
}

func TestCSRFMiddlewareForGET(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	e := echo.New()
	e.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())
	e.GET("/test", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"message": "GET success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "GET success")
}

func TestCSRFMiddlewareForCookie(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	e := echo.New()
	e.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())
	e.POST("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "POST success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "cookie_token"})
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "POST success")
}

func TestCSRFMiddlewareSuccess(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	e := echo.New()
	e.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())
	e.POST("/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "POST success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("X-CSRF-Token", "valid_token")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "POST success")
}
