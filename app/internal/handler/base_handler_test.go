package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/consts"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUuid(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedUuid := "test-uuid-1234"
	c.Set(consts.ContextKeys.Uuid, expectedUuid)

	handler := &BaseHandler{}
	uuid := handler.GetUuid(c)

	assert.Equal(t, expectedUuid, uuid)
}

func TestGetEmail(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedEmail := "test@example.com"
	c.Set(consts.ContextKeys.Email, expectedEmail)

	handler := &BaseHandler{}
	email := handler.GetEmail(c)

	assert.Equal(t, expectedEmail, email)
}

func TestGetRoomModel(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	room := model.Room{
		ID: primitive.NewObjectID(),
	}
	c.Set(consts.ContextKeys.RoomModel, room)

	handler := &BaseHandler{}
	gotRoom := handler.GetRoomModel(c)

	assert.Equal(t, room, gotRoom)
}

type ValidateHandler struct {
	BaseHandler
}

type testRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (h *ValidateHandler) ValidateRequestMockHandler(c echo.Context) error {
	var req testRequest
	if err := h.validateRequest(c, &req); err != nil {
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "valid request",
	})

}

func TestValidateRequestBindFail(t *testing.T) {
	e := echo.New()
	e.Validator = &usecase.CustomValidator{Validator: validator.New()}

	// Valid request
	body := `{ "name": 123 }`
	req := httptest.NewRequest(http.MethodPost, "/validate", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &ValidateHandler{}
	err := handler.ValidateRequestMockHandler(c)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "error")
}
