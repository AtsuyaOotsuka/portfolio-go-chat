package handler

import (
	"fmt"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/dto"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabapi"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
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
	mongoRoomSvc mongo_svc.RoomSvcInterface
	roomSvc      service.RoomSvcInterface
	dto          dto.RoomDtoInterface
}

func NewRoomHandler(
	mongoRoomSvc mongo_svc.RoomSvcInterface,
	roomSvc service.RoomSvcInterface,
	dto dto.RoomDtoInterface,
) *RoomHandler {
	return &RoomHandler{
		mongoRoomSvc: mongoRoomSvc,
		roomSvc:      roomSvc,
		dto:          dto,
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

	rooms, err := h.mongoRoomSvc.GetRoomList(uuid, target, ctx)
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

	InsertedID, err := h.mongoRoomSvc.CreateRoom(room, ctx)
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

	err = h.mongoRoomSvc.JoinRoom(req.RoomID, h.GetUuid(c), ctx)
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

	ctx := atylabapi.NewApiCtxSvc()
	defer ctx.Cancel()

	members, err := h.roomSvc.GetMemberInfos(room, ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"members": members,
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

	err := h.mongoRoomSvc.LeaveRoom(roomID, uuid, ctx)
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
	if !h.IsAdmin(c) {
		return c.JSON(400, echo.Map{
			"error": "Only admin can delete the room",
		})
	}

	roomID := c.Param("room_id")

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	err := h.mongoRoomSvc.DeleteRoom(roomID, ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "room deleted",
	})
}

type AddMemberRequest struct {
	MemberID string `json:"member_id" form:"member_id" validate:"required"`
}

func (h *RoomHandler) AddMember(c echo.Context) error {
	if !h.IsAdmin(c) {
		return c.JSON(400, echo.Map{
			"error": "Only admin can add members",
		})
	}

	var req AddMemberRequest
	if err := h.validateRequest(c, &req); err != nil {
		fmt.Println("Validation error:", err)
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	roomID := c.Param("room_id")

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	err := h.mongoRoomSvc.JoinRoom(roomID, req.MemberID, ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "member added",
	})
}

type RemoveMemberRequest struct {
	MemberID string `json:"member_id" form:"member_id" validate:"required"`
}

func (h *RoomHandler) RemoveMember(c echo.Context) error {
	if !h.IsAdmin(c) {
		return c.JSON(400, echo.Map{
			"error": "Only admin can remove members",
		})
	}

	var req RemoveMemberRequest
	if err := h.validateRequest(c, &req); err != nil {
		fmt.Println("Validation error:", err)
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	roomID := c.Param("room_id")

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	err := h.mongoRoomSvc.LeaveRoom(roomID, req.MemberID, ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"message": "member removed",
	})
}
