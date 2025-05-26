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

func TestGetTodoUseCase_Execute(t *testing.T) {
	t.Run("successfully retrieves todo", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		expectedTodo := todo.ReconstructTodo(
			todoID.UUID(),
			"Test Todo",
			"Test Body",
			todo.TodoStatusNotStarted,
			time.Now(),
			time.Now(),
		)

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Expect repository FindByID to be called within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				// Repository is called within transaction
				mockRepo.EXPECT().FindByID(ctx, todoID).Return(expectedTodo, nil)
				return fn(ctx)
			})

		useCase := &getTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		req := GetTodoRequest{ID: todoID}

		// When
		result, err := useCase.Execute(ctx, req)

		// Then
		require.NoError(t, err)
		require.Equal(t, expectedTodo, result)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		repoError := errors.New("repository error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Repository error occurs within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindByID(ctx, todoID).Return(nil, repoError)
				return fn(ctx)
			})

		useCase := &getTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		req := GetTodoRequest{ID: todoID}

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
		todoID := todo.TodoID(uuid.New())
		txError := errors.New("transaction error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Transaction itself fails
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(txError)

		useCase := &getTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		req := GetTodoRequest{ID: todoID}

		// When
		result, err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Nil(t, result)
		require.Contains(t, err.Error(), "failed to execute transaction")
	})
}
