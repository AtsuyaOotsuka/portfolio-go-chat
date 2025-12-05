package svc_mock

import (
	"github.com/stretchr/testify/mock"
)

type CsrfSvcMockStruct struct {
	mock.Mock
}

func (m *CsrfSvcMockStruct) CreateCSRFToken(
	timestamp int64,
	secret string,
) string {
	args := m.Called(timestamp, secret)
	return args.String(0)
}

func (m *CsrfSvcMockStruct) Verify(
	token string,
	secret string,
	timestamp int64,
) error {
	args := m.Called(token, secret, timestamp)
	return args.Error(0)
}
