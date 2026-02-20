package testutil

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)

	if args.Get(0) == nil {
		return fn(ctx)
	}
	return args.Error(0)
}

func NewMockTransactionManager() *MockTransactionManager {
	m := &MockTransactionManager{}
	m.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
	return m
}

func NewFailingTransactionManager(err error) *MockTransactionManager {
	m := &MockTransactionManager{}
	m.On("WithTransaction", mock.Anything, mock.Anything).Return(err)
	return m
}

type PassthroughTransactionManager struct{}

func (m *PassthroughTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func NewPassthroughTransactionManager() *PassthroughTransactionManager {
	return &PassthroughTransactionManager{}
}
