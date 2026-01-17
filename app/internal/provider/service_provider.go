package provider

import (
	"os"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/service/mongo_svc"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabapi"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
)

func (p *Provider) bindMongoRoomSvc() mongo_svc.RoomSvcInterface {
	return mongo_svc.NewRoomSvcStruct(
		p.bindMongoSvc(),
	)
}

func (p *Provider) bindMongoMessageSvc() mongo_svc.MessageSvcInterface {
	return mongo_svc.NewMessageSvcStruct(
		p.bindMongoSvc(),
	)
}

func (p *Provider) bindCsrfSvc() service.CsrfSvcInterface {
	return service.NewCsrfSvcStruct(
		atylabcsrf.NewCsrfPkgStruct(),
	)
}

func (p *Provider) bindRoomSvc() service.RoomSvcInterface {
	return service.NewRoomSvc(
		p.bindRedisSvc(),
		p.bindMongoRoomSvc(),
		atylabapi.NewApiPostStruct(
			os.Getenv("COMMON_KEY"),
			os.Getenv("API_BASE_URL"),
		),
	)
}
