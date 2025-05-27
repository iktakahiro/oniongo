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

func TestStartTodoUseCase_Execute(t *testing.T) {
	t.Run("successfully starts todo", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		req := StartTodoRequest{ID: todoID}

		existingTodo := todo.ReconstructTodo(
			todoID.UUID(),
			"Test Todo",
			"Test Body",
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

		useCase := &startTodoUseCase{
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
		req := StartTodoRequest{ID: todoID}
		findError := errors.New("todo not found")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// FindByID error occurs within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindByID(ctx, todoID).Return(nil, findError)
				return fn(ctx)
			})

		useCase := &startTodoUseCase{
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
		req := StartTodoRequest{ID: todoID}

		existingTodo := todo.ReconstructTodo(
			todoID.UUID(),
			"Test Todo",
			"Test Body",
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

		useCase := &startTodoUseCase{
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
		req := StartTodoRequest{ID: todoID}
		txError := errors.New("transaction error")

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Transaction itself fails
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(txError)

		useCase := &startTodoUseCase{
			todoRepository: mockRepo,
			txRunner:       mockTxRunner,
		}

		// When
		err := useCase.Execute(ctx, req)

		// Then
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to execute transaction")
	})

	t.Run("returns error when todo cannot be started", func(t *testing.T) {
		// Given
		ctx := context.Background()
		todoID := todo.TodoID(uuid.New())
		req := StartTodoRequest{ID: todoID}

		// Todo that is already completed cannot be started
		existingTodo := todo.ReconstructTodo(
			todoID.UUID(),
			"Test Todo",
			"Test Body",
			todo.TodoStatusCompleted,
			time.Now(),
			time.Now(),
		)

		mockRepo := mock_todo.NewMockTodoRepository(t)
		mockTxRunner := mock_uow.NewMockTransactionRunner(t)

		// Start error occurs within transaction
		mockTxRunner.EXPECT().RunInTx(ctx, mock.AnythingOfType("func(context.Context) error")).
			RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
				mockRepo.EXPECT().FindByID(ctx, todoID).Return(existingTodo, nil)
				return fn(ctx)
			})

		useCase := &startTodoUseCase{
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