package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/svc_mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRoomMiddlewareHandler(t *testing.T) {

	room := model.Room{
		ID:      primitive.NewObjectID(),
		OwnerID: "test-uuid-1234",
		Members: []string{"test-uuid-1234"},
	}

	mockRoomSvc := new(svc_mock.RoomSvcMock)
	mockRoomSvc.On("GetRoom", room.ID.Hex(), mock.Anything).Return(room, nil)
	mockRoomSvc.On("IsOwner", room, "test-uuid-1234").Return(true)
	mockRoomSvc.On("IsMember", room, "test-uuid-1234").Return(true)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("uuid", "test-uuid-1234")
			return next(c)
		}
	})
	e.Use(NewRoomMiddleware(mockRoomSvc).Handler())
	e.GET("/test/:room_id", func(c echo.Context) error {
		roomModel := c.Get("room_model").(model.Room)
		isAdmin := c.Get("is_admin").(bool)
		isMember := c.Get("is_member").(bool)

		assert.True(t, isAdmin)
		assert.True(t, isMember)
		assert.Equal(t, room.ID, roomModel.ID)

		return c.JSON(200, echo.Map{"message": "success"})
	})
	req := httptest.NewRequest(http.MethodGet, "/test/"+room.ID.Hex(), nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	c := e.NewContext(req, w)
	c.SetParamNames("room_id")
	c.SetParamValues(room.ID.Hex())

	assert.Equal(t, http.StatusOK, w.Code)
	mockRoomSvc.AssertExpectations(t)
}

func TestRoomMiddlewareHandlerNotFound(t *testing.T) {

	room := model.Room{}

	mockRoomSvc := new(svc_mock.RoomSvcMock)
	mockRoomSvc.On("GetRoom", room.ID.Hex(), mock.Anything).Return(room, assert.AnError)

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("uuid", "test-uuid-1234")
			return next(c)
		}
	})
	e.Use(NewRoomMiddleware(mockRoomSvc).Handler())
	e.GET("/test/:room_id", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"message": "success"})
	})
	req := httptest.NewRequest(http.MethodGet, "/test/"+room.ID.Hex(), nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	c := e.NewContext(req, w)
	c.SetParamNames("room_id")
	c.SetParamValues(room.ID.Hex())

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockRoomSvc.AssertExpectations(t)
}
