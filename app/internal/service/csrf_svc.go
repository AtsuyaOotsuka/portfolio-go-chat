package service

import "github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"

type CsrfSvcInterface interface {
	CreateCSRFToken(timestamp int64, secret string) string
	Verify(token string, secret string, timestamp int64) error
}

type CsrfSvcStruct struct {
	csrf atylabcsrf.CsrfPkgInterface
}

func NewCsrfSvcStruct(
	csrf atylabcsrf.CsrfPkgInterface,
) *CsrfSvcStruct {
	return &CsrfSvcStruct{
		csrf: csrf,
	}
}

func (s *CsrfSvcStruct) CreateCSRFToken(
	timestamp int64,
	secret string,
) string {
	nonceStr := s.csrf.GenerateNonceString()
	return s.csrf.GenerateCSRFCookieToken(secret, timestamp, nonceStr)
}

func (s *CsrfSvcStruct) Verify(
	token string,
	secret string,
	timestamp int64,
) error {
	return s.csrf.ValidateCSRFCookieToken(token, secret, timestamp)
}
