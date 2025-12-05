package provider

import "testing"

func TestBindHealthCheckHandler(t *testing.T) {
	provider := NewProvider()
	healthCheckHandler := provider.BindHealthCheckHandler()

	if healthCheckHandler == nil {
		t.Fatal("BindHealthCheckHandler returned nil")
	}
}
