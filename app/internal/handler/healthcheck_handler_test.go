package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler_Check(t *testing.T) {
	var err error
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("uuid", "test-uuid-1234")
	c.Set("email", "test@example.com")

	handler := NewHealthCheckHandler()
	err = handler.Check(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	result := map[string]string{}
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
	assert.Equal(t, "test-uuid-1234", result["uuid"])
	assert.Equal(t, "test@example.com", result["email"])
}
