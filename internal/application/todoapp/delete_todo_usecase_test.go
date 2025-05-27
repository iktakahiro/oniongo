package todoapp

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
	"github.com/iktakahiro/oniongo/internal/mocks/application/mock_uow"
	"github.com/iktakahiro/oniongo/internal/mocks/domain/mock_todo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteTodoUseCase_Execute(t *testing.T) {
	t.Run("successfully deletes todo", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		req := DeleteTodoRequest{ID: todoID}

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Expect repository Delete to be called within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().Delete(ctx, todoID).Return(nil)
				return fn(ctx)
			})

		useCase := &deleteTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		err := useCase.Execute(ctx, req)

		// Then
		require.NoError(t, err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		req := DeleteTodoRequest{ID: todoID}
		repoError := errors.New("repository error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Repository error occurs within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().Delete(ctx, todoID).Return(repoError)
				return fn(ctx)
			})

		useCase := &deleteTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to execute transaction")
	})

	t.Run("returns error when transaction fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		req := DeleteTodoRequest{ID: todoID}
		txError := errors.New("transaction error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Transaction itself fails
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(txError)

		useCase := &deleteTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to execute transaction")
	})
} 