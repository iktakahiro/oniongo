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

func TestUpdateTodoUseCase_Execute(t *testing.T) {
	t.Run("successfully updates todo", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		req := UpdateTodoRequest{
			ID:    todoID,
			Title: "Updated Title",
			Body:  "Updated Body",
		}

		existingTodo := todo.ReconstructTodo(
			todoID.UUID(),
			"Original Title",
			"Original Body",
			todo.TodoStatusNotStarted,
			time.Now(),
			time.Now(),
		)

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Expect repository operations to be called within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindByID(ctx, todoID).Return(existingTodo, nil)
				mockRepo.EXPECT().Update(ctx, existingTodo).Return(nil)
				return fn(ctx)
			})

		useCase := &updateTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		err := useCase.Execute(ctx, req)

		// Then
		require.NoError(t, err)
	})

	t.Run("returns error when todo not found", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		req := UpdateTodoRequest{
			ID:    todoID,
			Title: "Updated Title",
			Body:  "Updated Body",
		}
		findError := errors.New("todo not found")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// FindByID error occurs within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindByID(ctx, todoID).Return(nil, findError)
				return fn(ctx)
			})

		useCase := &updateTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to execute transaction")
	})

	t.Run("returns error when update repository fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		req := UpdateTodoRequest{
			ID:    todoID,
			Title: "Updated Title",
			Body:  "Updated Body",
		}

		existingTodo := todo.ReconstructTodo(
			todoID.UUID(),
			"Original Title",
			"Original Body",
			todo.TodoStatusNotStarted,
			time.Now(),
			time.Now(),
		)
		updateError := errors.New("update failed")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Update error occurs within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindByID(ctx, todoID).Return(existingTodo, nil)
				mockRepo.EXPECT().Update(ctx, existingTodo).Return(updateError)
				return fn(ctx)
			})

		useCase := &updateTodoUseCase{
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
		req := UpdateTodoRequest{
			ID:    todoID,
			Title: "Updated Title",
			Body:  "Updated Body",
		}
		txError := errors.New("transaction error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Transaction itself fails
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(txError)

		useCase := &updateTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to execute transaction")
	})

	t.Run("returns error when setting invalid title", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		req := UpdateTodoRequest{
			ID:    todoID,
			Title: "", // Empty title should cause domain error
			Body:  "Updated Body",
		}

		existingTodo := todo.ReconstructTodo(
			todoID.UUID(),
			"Original Title",
			"Original Body",
			todo.TodoStatusNotStarted,
			time.Now(),
			time.Now(),
		)

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// SetTitle error occurs within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindByID(ctx, todoID).Return(existingTodo, nil)
				return fn(ctx)
			})

		useCase := &updateTodoUseCase{
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