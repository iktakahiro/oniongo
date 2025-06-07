package todoapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application/uow"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
	"github.com/samber/do"
)

type CompleteTodoRequest struct {
	ID todo.TodoID
}

// CompleteTodoUseCase is the interface that wraps the basic CompleteTodo operation.
type CompleteTodoUseCase interface {
	Execute(ctx context.Context, req CompleteTodoRequest) error
}

// completeTodoUseCase is the implementation of the CompleteTodoUseCase interface.
type completeTodoUseCase struct {
	todoRepository todo.TodoRepository
	txRunner       uow.TransactionRunner
}

// NewCompleteTodoUseCase creates a new CompleteTodoUseCase.
func NewCompleteTodoUseCase(i *do.Injector) (CompleteTodoUseCase, error) {
	todoRepository, err := do.Invoke[todo.TodoRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke todo repository: %w", err)
	}
	txRunner, err := do.Invoke[uow.TransactionRunner](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke transaction manager: %w", err)
	}

	return &completeTodoUseCase{
		todoRepository: todoRepository,
		txRunner:       txRunner,
	}, nil
}

// Execute completes a Todo by changing its status to completed.
func (u *completeTodoUseCase) Execute(ctx context.Context, req CompleteTodoRequest) error {
	err := u.txRunner.RunInTx(ctx, func(ctx context.Context) error {
		foundTodo, err := u.todoRepository.FindByID(ctx, req.ID)
		if err != nil {
			var notFoundErr *todo.NotFoundError
			if errors.As(err, &notFoundErr) {
				return err
			}
			return fmt.Errorf("failed to find todo: %w", err)
		}

		if err := foundTodo.Complete(); err != nil {
			// Preserve domain errors
			var stateErr *todo.StateError
			if errors.As(err, &stateErr) {
				return err
			}
			return fmt.Errorf("failed to complete todo: %w", err)
		}

		if err := u.todoRepository.Update(ctx, foundTodo); err != nil {
			return fmt.Errorf("failed to update todo: %w", err)
		}
		return nil
	})
	if err != nil {
		// Preserve domain errors
		var notFoundErr *todo.NotFoundError
		var stateErr *todo.StateError
		if errors.As(err, &notFoundErr) || errors.As(err, &stateErr) {
			return err
		}
		return fmt.Errorf("failed to execute transaction: %w", err)
	}
	return nil
}
