package testutil

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockTransactionManager_ExecutesFunction(t *testing.T) {
	txManager := NewMockTransactionManager()
	ctx := context.Background()

	executed := false
	err := txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		executed = true
		return nil
	})

	require.NoError(t, err)
	assert.True(t, executed, "Transaction function should be executed")
	txManager.AssertExpectations(t)
}

func TestMockTransactionManager_PropagatesError(t *testing.T) {
	txManager := NewMockTransactionManager()
	ctx := context.Background()

	expectedErr := errors.New("business logic error")
	err := txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return expectedErr
	})

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestFailingTransactionManager_ReturnsError(t *testing.T) {
	expectedErr := errors.New("transaction failed")
	txManager := NewFailingTransactionManager(expectedErr)
	ctx := context.Background()

	executed := false
	err := txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		executed = true
		return nil
	})

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.False(t, executed, "Transaction function should NOT be executed on failure")
}

func TestPassthroughTransactionManager_ExecutesFunction(t *testing.T) {
	txManager := NewPassthroughTransactionManager()
	ctx := context.Background()

	executed := false
	err := txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		executed = true
		assert.Equal(t, ctx, txCtx, "Context should be passed through")
		return nil
	})

	require.NoError(t, err)
	assert.True(t, executed, "Transaction function should be executed")
}

func TestPassthroughTransactionManager_PropagatesError(t *testing.T) {
	txManager := NewPassthroughTransactionManager()
	ctx := context.Background()

	expectedErr := errors.New("some error")
	err := txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		return expectedErr
	})

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
