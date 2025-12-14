package provider

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
)

func TestBindHealthCheckHandler(t *testing.T) {
	provider := NewProvider(usecase.NewMongo())
	healthCheckHandler := provider.BindHealthCheckHandler()

	if healthCheckHandler == nil {
		t.Fatal("BindHealthCheckHandler returned nil")
	}
}

func TestBindRoomHandler(t *testing.T) {
	provider := NewProvider(usecase.NewMongo())
	roomHandler := provider.BindRoomHandler()

	if roomHandler == nil {
		t.Fatal("BindRoomHandler returned nil")
	}
}

func TestBindMessageHandler(t *testing.T) {
	provider := NewProvider(usecase.NewMongo())
	messageHandler := provider.BindMessageHandler()

	if messageHandler == nil {
		t.Fatal("BindMessageHandler returned nil")
	}
}
