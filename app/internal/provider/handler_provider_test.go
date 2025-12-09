package provider

import "testing"

func TestBindHealthCheckHandler(t *testing.T) {
	provider := NewProvider()
	healthCheckHandler := provider.BindHealthCheckHandler()

	if healthCheckHandler == nil {
		t.Fatal("BindHealthCheckHandler returned nil")
	}
}

func TestBindRoomHandler(t *testing.T) {
	provider := NewProvider()
	roomHandler := provider.BindRoomHandler()

	if roomHandler == nil {
		t.Fatal("BindRoomHandler returned nil")
	}
}

func TestBindMessageHandler(t *testing.T) {
	provider := NewProvider()
	messageHandler := provider.BindMessageHandler()

	if messageHandler == nil {
		t.Fatal("BindMessageHandler returned nil")
	}
}
