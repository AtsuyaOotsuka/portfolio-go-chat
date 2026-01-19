package provider

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
)

func TestBindBindCsrfMiddleware(t *testing.T) {
	provider := NewProvider(usecase.NewMongo(), usecase.NewRedis())
	csrfMiddleware := provider.BindCsrfMiddleware()

	if csrfMiddleware == nil {
		t.Fatal("BindCsrfMiddleware returned nil")
	}
}

func TestBindJwtMiddleware(t *testing.T) {
	provider := NewProvider(usecase.NewMongo(), usecase.NewRedis())
	jwtMiddleware := provider.BindJwtMiddleware()

	if jwtMiddleware == nil {
		t.Fatal("BindJwtMiddleware returned nil")
	}
}

func TestBindRoomMiddleware(t *testing.T) {
	provider := NewProvider(usecase.NewMongo(), usecase.NewRedis())
	roomMiddleware := provider.BindRoomMiddleware()

	if roomMiddleware == nil {
		t.Fatal("BindRoomMiddleware returned nil")
	}
}
