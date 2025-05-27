package todoapp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
	"github.com/iktakahiro/oniongo/internal/mocks/application/mock_uow"
	"github.com/iktakahiro/oniongo/internal/mocks/domain/mock_todo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetTodosUseCase_Execute(t *testing.T) {
	t.Run("successfully retrieves todos", func(t *testing.T) {
		// Given
		ctx := context.Background()
		req := GetTodosRequest{}

		expectedTodos := []*todo.Todo{
			todo.ReconstructTodo(
				uuid.New(),
				"Todo 1",
				"Body 1",
				todo.TodoStatusNotStarted,
				time.Now(),
				time.Now(),
			),
			todo.ReconstructTodo(
				uuid.New(),
				"Todo 2",
				"Body 2",
				todo.TodoStatusInProgress,
				time.Now(),
				time.Now(),
			),
		}

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Expect repository FindAll to be called within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindAll(ctx).Return(expectedTodos, nil)
				return fn(ctx)
			})

		useCase := &getTodosUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		result, err := useCase.Execute(ctx, req)

		// Then
		require.NoError(t, err)
		require.Equal(t, expectedTodos, result)
	})

	t.Run("successfully retrieves empty todos list", func(t *testing.T) {
		// Given
		ctx := context.Background()
		req := GetTodosRequest{}
		expectedTodos := []*todo.Todo{}

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Expect repository FindAll to be called within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindAll(ctx).Return(expectedTodos, nil)
				return fn(ctx)
			})

		useCase := &getTodosUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		result, err := useCase.Execute(ctx, req)

		// Then
		require.NoError(t, err)
		require.Equal(t, expectedTodos, result)
		require.Len(t, result, 0)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		req := GetTodosRequest{}
		repoError := errors.New("repository error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Repository error occurs within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindAll(ctx).Return(nil, repoError)
				return fn(ctx)
			})

		useCase := &getTodosUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		result, err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Nil(t, result)
		require.Contains(t, err.Error(), "failed to execute transaction")
	})

	t.Run("returns error when transaction fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		req := GetTodosRequest{}
		txError := errors.New("transaction error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Transaction itself fails
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(txError)

		useCase := &getTodosUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		result, err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Nil(t, result)
		require.Contains(t, err.Error(), "failed to execute transaction")
	})
} 