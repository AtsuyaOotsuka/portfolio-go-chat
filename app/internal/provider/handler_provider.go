package provider

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/dto"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/handler"
)

func (p *Provider) BindHealthCheckHandler() *handler.HealthCheckHandler {
	return handler.NewHealthCheckHandler()
}

func (p *Provider) BindRoomHandler() *handler.RoomHandler {
	return handler.NewRoomHandler(
		p.bindMongoRoomSvc(),
		dto.NewRoomDtoStruct(),
	)
}

func (p *Provider) BindMessageHandler() *handler.MessageHandler {
	return handler.NewMessageHandler(
		p.bindMongoMessageSvc(),
		dto.NewMessageDtoStruct(),
	)
}
