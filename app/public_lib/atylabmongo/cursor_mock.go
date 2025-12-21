package atylabmongo

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MongoCursorStructMock struct {
	mock.Mock
}

func (m *MongoCursorStructMock) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MongoCursorStructMock) Decode(val interface{}) error {
	args := m.Called(val)
	return args.Error(0)
}

func (m *MongoCursorStructMock) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MongoCursorStructMock) All(ctx context.Context, result interface{}) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}
