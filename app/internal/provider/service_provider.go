package provider

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc"
)

func (p *Provider) bindRoomSvc() *mongo_svc.RoomSvcStruct {
	return mongo_svc.NewRoomSvcStruct(
		p.bindMongoSvc(),
	)
}
