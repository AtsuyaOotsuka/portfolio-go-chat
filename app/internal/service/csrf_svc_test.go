package service

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
)

func TestNewCsrfSvcStruct(t *testing.T) {
	svc := NewCsrfSvcStruct(&atylabcsrf.CsrfPkgMockStruct{})
	if svc == nil {
		t.Error("expected non-nil CsrfSvcStruct")
	}
}

func TestCreateCSRFToken(t *testing.T) {
	csrfMock := &atylabcsrf.CsrfPkgMockStruct{}

	cvs := CsrfSvcStruct{
		csrf: csrfMock,
	}

	token := cvs.CreateCSRFToken(1234567890, "test_secret")

	if token != "mocked_csrf_token" {
		t.Errorf("expected 'mocked_csrf_token', got '%s'", token)
	}
}

func TestVerify(t *testing.T) {
	csrfMock := &atylabcsrf.CsrfPkgMockStruct{}

	cvs := CsrfSvcStruct{
		csrf: csrfMock,
	}

	err := cvs.Verify("test_token", "test_secret", 1234567890)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
