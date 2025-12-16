package handler

import (
	"fmt"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/dto"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"github.com/labstack/echo/v4"
)

type RoomHandlerInterface interface {
	List(c echo.Context) error
	Create(c echo.Context) error
	Join(c echo.Context) error
	Members(c echo.Context) error
	Leave(c echo.Context) error
	Delete(c echo.Context) error
	AddMember(c echo.Context) error
	RemoveMember(c echo.Context) error
}

type RoomHandler struct {
	BaseHandler
	service mongo_svc.RoomSvcInterface
	dto     dto.RoomDtoInterface
}

func NewRoomHandler(
	service mongo_svc.RoomSvcInterface,
	dto dto.RoomDtoInterface,
) *RoomHandler {
	return &RoomHandler{
		service: service,
		dto:     dto,
	}
}

func (h *RoomHandler) List(c echo.Context) error {
	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	target := c.QueryParam("target")
	if target == "" {
		target = "all"
	}
	uuid := h.GetUuid(c)

	rooms, err := h.service.GetRoomList(uuid, target, ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"rooms": h.dto.ResponseRoomList(rooms, uuid),
	})
}

type CreateRoomRequest struct {
	Name      string `json:"name" form:"name" validate:"required"`
	IsPrivate bool   `json:"is_private" form:"is_private"`
}

func (h *RoomHandler) Create(c echo.Context) error {
	var req CreateRoomRequest
	if err := h.validateRequest(c, &req); err != nil {
		fmt.Println("Validation error:", err)
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	uuid := h.GetUuid(c)

	room := model.Room{
		Name:      req.Name,
		OwnerID:   uuid,
		CreatedAt: time.Now(),
		Members:   []string{uuid},
		IsPrivate: req.IsPrivate,
	}

	InsertedID, err := h.service.CreateRoom(room, ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}
	fmt.Print(InsertedID)

	// Implement room creation logic here
	return c.JSON(200, echo.Map{
		"message":    "Room created successfully",
		"room_id":    InsertedID,
		"room_name":  room.Name,
		"created_at": room.CreatedAt,
	})
}

type RoomJoinRequest struct {
	RoomID string `json:"room_id" form:"room_id" validate:"required"`
}

func (h *RoomHandler) Join(c echo.Context) error {
	var req RoomJoinRequest
	var err error
	if err := h.validateRequest(c, &req); err != nil {
		fmt.Println("Validation error:", err)
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	if h.IsMember(c) {
		return c.JSON(400, echo.Map{
			"error": "Already a member of the room",
		})
	}

	err = h.service.JoinRoom(req.RoomID, h.GetUuid(c), ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "Joined room successfully",
	})
}

func (h *RoomHandler) Members(c echo.Context) error {
	if !h.IsMember(c) {
		return c.JSON(400, echo.Map{
			"error": "Not a member of the room",
		})
	}

	room := h.GetRoomModel(c)

	return c.JSON(200, echo.Map{
		"members": room.Members,
	})
}

func (h *RoomHandler) Leave(c echo.Context) error {

	if !h.IsMember(c) {
		return c.JSON(400, echo.Map{
			"error": "Not a member of the room",
		})
	}

	if h.IsAdmin(c) {
		return c.JSON(400, echo.Map{
			"error": "Admin cannot leave the room",
		})
	}

	uuid := h.GetUuid(c)
	roomID := c.Param("room_id")

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	err := h.service.LeaveRoom(roomID, uuid, ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "left room",
	})
}

func (h *RoomHandler) Delete(c echo.Context) error {
	// Implement room deletion logic here
	return c.JSON(200, echo.Map{
		"message": "room deleted",
	})
}

func (h *RoomHandler) AddMember(c echo.Context) error {
	// Implement add member logic here
	return c.JSON(200, echo.Map{
		"message": "member added",
	})
}

func (h *RoomHandler) RemoveMember(c echo.Context) error {
	// Implement remove member logic here
	return c.JSON(200, echo.Map{
		"message": "member removed",
	})
}
