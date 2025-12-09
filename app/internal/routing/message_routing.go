package routing

import "github.com/AtsuyaOotsuka/portfolio-go-chat/internal/handler"

func (r *Routing) MessageRoute(
	handler handler.MessageHandlerInterface,
) {
	messageGroup := r.echo.Group("/message")

	messageGroup.GET("/:room_id/list", handler.List)
	messageGroup.POST("/:room_id/send", handler.Send)
	messageGroup.POST("/:room_id/read", handler.Read)
	messageGroup.DELETE("/:room_id/delete", handler.Delete)

	r.Finalize(messageGroup)
}
