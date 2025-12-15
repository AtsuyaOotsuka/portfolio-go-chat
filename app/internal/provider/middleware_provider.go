package provider

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
)

func (p *Provider) BindCsrfMiddleware() middleware.CSRFMiddlewareInterface {
	return middleware.NewCSRFMiddleware(
		p.bindCsrfSvc(),
	)
}

func (p *Provider) BindJwtMiddleware() middleware.JWTMiddlewareInterface {
	return middleware.NewJWTMiddleware(
		atylabjwt.NewJwtSvc(),
	)
}

func (p *Provider) BindRoomMiddleware() middleware.RoomMVMiddlewareInterface {
	return middleware.NewRoomMiddleware(
		p.bindRoomSvc(),
	)
}
