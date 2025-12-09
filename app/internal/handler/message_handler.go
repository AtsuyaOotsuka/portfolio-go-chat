package handler

import (
	"fmt"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/dto"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/model"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc"
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
	roomSvc    mongo_svc.RoomSvcInterface
	dto        dto.MessageDtoInterface
}

func NewMessageHandler(
	messageSvc mongo_svc.MessageSvcInterface,
	roomSvc mongo_svc.RoomSvcInterface,
	dto dto.MessageDtoInterface,
) *MessageHandler {
	return &MessageHandler{
		messageSvc: messageSvc,
		roomSvc:    roomSvc,
		dto:        dto,
	}
}

func (h *MessageHandler) List(c echo.Context) error {
	roomID := c.Param("room_id")
	uuid := h.GetUuid(c)

	if err := h.roomSvc.IsJoinedRoom(roomID, uuid); err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	messages, err := h.messageSvc.GetMessageList(roomID)
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

	roomID := c.Param("room_id")
	uuid := h.GetUuid(c)

	if err := h.roomSvc.IsJoinedRoom(roomID, uuid); err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}

	message := model.Message{
		RoomID:        roomID,
		Sender:        uuid,
		Message:       req.Message,
		CreatedAt:     time.Now(),
		IsReadUserIds: []string{uuid},
	}

	messageId, err := h.messageSvc.SendMessage(message)
	if err != nil {
		return c.JSON(500, echo.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(200, echo.Map{
		"message_id": messageId,
	})
}

func (h *MessageHandler) Read(c echo.Context) error {
	return nil
}

func (h *MessageHandler) Delete(c echo.Context) error {
	return nil
}
