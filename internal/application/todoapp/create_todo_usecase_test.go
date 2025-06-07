package todoapp

import (
	"context"
	"errors"
	"testing"

	"github.com/iktakahiro/oniongo/internal/mocks/application/mock_uow"
	"github.com/iktakahiro/oniongo/internal/mocks/domain/mock_todo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateTodoUseCase_Execute(t *testing.T) {
	t.Run("successfully creates todo", func(t *testing.T) {
		// Given
		ctx := context.Background()
		req := CreateTodoRequest{
			Title: "Test Todo",
			Body:  "Test Body",
		}

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Expect repository Create to be called within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				// Repository is called within transaction
				mockRepo.EXPECT().Create(ctx, mock.AnythingOfType("*todo.Todo")).Return(nil)
				return fn(ctx)
			})

		useCase := &createTodoUseCase{
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
		req := CreateTodoRequest{
			Title: "Test Todo",
			Body:  "Test Body",
		}
		repoError := errors.New("repository error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Repository error occurs within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().Create(ctx, mock.AnythingOfType("*todo.Todo")).Return(repoError)
				return fn(ctx)
			})

		useCase := &createTodoUseCase{
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
		req := CreateTodoRequest{
			Title: "Test Todo",
			Body:  "Test Body",
		}
		txError := errors.New("transaction error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Transaction itself fails
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(txError)

		useCase := &createTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to execute transaction")
	})

	t.Run("returns error when todo creation fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		req := CreateTodoRequest{
			Title: "", // Empty title should cause domain error
			Body:  "Test Body",
		}

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		useCase := &createTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Contains(t, err.Error(), "title is required")
		// Transaction should not be called when todo creation fails
		mockTxRunner.AssertNotCalled(t, "RunInTx")
	})
}
