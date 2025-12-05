package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestJwtHandler(t *testing.T) {
	funcs.WithEnv("JWT_SECRET_KEY", "testsecretkey", t, func() {
		jwt := "tokenstring"
		mockJwt := new(atylabjwt.JwtMock)
		mockJwt.On("Validate", "testsecretkey", jwt).Return(nil)
		mockJwt.On("GetUUID").Return("test-uuid")
		mockJwt.On("GetEmail").Return("test@example.com")

		e := echo.New()
		e.Use(NewJWTMiddleware(mockJwt).Handler())

		e.GET("/test", func(c echo.Context) error {
			return c.JSON(200, echo.Map{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+jwt)
		w := httptest.NewRecorder()

		e.GET("/test", func(c echo.Context) error {
			userID := mockJwt.GetUUID()
			email := mockJwt.GetEmail()

			assert.Equal(t, userID, c.Get("uuid"))
			assert.Equal(t, email, c.Get("email"))

			return c.JSON(200, echo.Map{"uuid": userID, "email": email})
		})
		e.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		result := map[string]string{}
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, "test-uuid", result["uuid"])
		assert.Equal(t, "test@example.com", result["email"])
	})
}

func TestJwtHandler_NoToken(t *testing.T) {
	mockJwt := new(atylabjwt.JwtMock)

	e := echo.New()
	e.Use(NewJWTMiddleware(mockJwt).Handler())

	e.GET("/test", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "not set jwt token")
}

func TestJwtHandler_InvalidToken(t *testing.T) {
	funcs.WithEnv("JWT_SECRET_KEY", "testsecretkey", t, func() {
		invalidJwt := "invalidtokenstring"
		mockJwt := new(atylabjwt.JwtMock)
		mockJwt.On("Validate", "testsecretkey", invalidJwt).Return(assert.AnError)

		e := echo.New()
		e.Use(NewJWTMiddleware(mockJwt).Handler())

		e.GET("/test", func(c echo.Context) error {
			return c.JSON(200, echo.Map{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+invalidJwt)
		w := httptest.NewRecorder()
		e.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), assert.AnError.Error())
	})
}
