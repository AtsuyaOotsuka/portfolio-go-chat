package provider

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc"
)

func (p *Provider) bindRoomSvc() *mongo_svc.RoomSvcStruct {
	return mongo_svc.NewRoomSvcStruct(
		p.bindMongoSvc(),
	)
}

func (p *Provider) bindMessageSvc() *mongo_svc.MessageSvcStruct {
	return mongo_svc.NewMessageSvcStruct(
		p.bindMongoSvc(),
	)
}
