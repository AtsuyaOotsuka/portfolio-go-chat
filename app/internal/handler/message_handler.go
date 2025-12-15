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

type MessageHandlerInterface interface {
	List(c echo.Context) error
	Send(c echo.Context) error
	Read(c echo.Context) error
	Delete(c echo.Context) error
}

type MessageHandler struct {
	BaseHandler
	messageSvc mongo_svc.MessageSvcInterface
	dto        dto.MessageDtoInterface
}

func NewMessageHandler(
	messageSvc mongo_svc.MessageSvcInterface,
	dto dto.MessageDtoInterface,
) *MessageHandler {
	return &MessageHandler{
		messageSvc: messageSvc,
		dto:        dto,
	}
}

func (h *MessageHandler) List(c echo.Context) error {
	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	roomID := c.Param("room_id")
	uuid := h.GetUuid(c)
	if !h.IsMember(c) {
		return c.JSON(403, echo.Map{
			"error": "You are not a member of this room.",
		})
	}

	messages, err := h.messageSvc.GetMessageList(roomID, ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"messages": h.dto.ResponseMessageList(messages, uuid),
	})
}

type SendMessageRequest struct {
	Message string `json:"message" form:"message" validate:"required"`
}

func (h *MessageHandler) Send(c echo.Context) error {
	var req SendMessageRequest
	if err := h.validateRequest(c, &req); err != nil {
		fmt.Println("Validation error:", err)
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	roomID := c.Param("room_id")
	uuid := h.GetUuid(c)

	if !h.IsMember(c) {
		return c.JSON(403, echo.Map{
			"error": "You are not a member of this room.",
		})
	}

	message := model.Message{
		RoomID:        roomID,
		Sender:        uuid,
		Message:       req.Message,
		CreatedAt:     time.Now(),
		IsReadUserIds: []string{uuid},
	}

	messageId, err := h.messageSvc.SendMessage(message, ctx)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"message_id": messageId,
	})
}

type ReadMessageRequest struct {
	MessageIds []string `json:"message_ids" form:"message_ids" validate:"required"`
}

func (h *MessageHandler) Read(c echo.Context) error {
	var req ReadMessageRequest
	if err := h.validateRequest(c, &req); err != nil {
		fmt.Println("Validation error:", err)
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	roomID := c.Param("room_id")
	uuid := h.GetUuid(c)
	if !h.IsMember(c) {
		return c.JSON(403, echo.Map{
			"error": "You are not a member of this room.",
		})
	}
	messageIDs := req.MessageIds

	fmt.Println("Marking messages as read:", messageIDs, "for user:", uuid, "in room:", roomID)
	if err := h.messageSvc.ReadMessages(messageIDs, roomID, uuid, ctx); err != nil {
		fmt.Println("Failed to mark messages as read:", err)
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status": "success",
	})
}

type DeleteMessageRequest struct {
	MessageId string `json:"message_id" form:"message_id" validate:"required"`
}

func (h *MessageHandler) Delete(c echo.Context) error {
	var req DeleteMessageRequest
	if err := h.validateRequest(c, &req); err != nil {
		fmt.Println("Validation error:", err)
		return c.JSON(400, echo.Map{
			"error": err.Error(),
		})
	}

	ctx := atylabmongo.NewMongoCtxSvc()
	defer ctx.Cancel()

	roomID := c.Param("room_id")
	uuid := h.GetUuid(c)
	if !h.IsMember(c) {
		return c.JSON(403, echo.Map{
			"error": "You are not a member of this room.",
		})
	}

	messageID := req.MessageId

	if err := h.messageSvc.IsSender(messageID, roomID, uuid, ctx); err != nil {
		if !h.IsAdmin(c) {
			return c.JSON(403, echo.Map{
				"error": "You are not authorized to delete this message.",
			})
		}
	}

	if err := h.messageSvc.DeleteMessage(messageID, roomID, ctx); err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(200, echo.Map{
		"status": "success",
	})
}
