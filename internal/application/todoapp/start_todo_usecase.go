package todoapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application/uow"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
	"github.com/samber/do"
)

type StartTodoRequest struct {
	ID todo.TodoID
}

// StartTodoUseCase is the interface that wraps the basic StartTodo operation.
type StartTodoUseCase interface {
	Execute(ctx context.Context, req StartTodoRequest) error
}

// startTodoUseCase is the implementation of the StartTodoUseCase interface.
type startTodoUseCase struct {
	todoRepository todo.TodoRepository
	txRunner       uow.TransactionRunner
}

// NewStartTodoUseCase creates a new StartTodoUseCase.
func NewStartTodoUseCase(i *do.Injector) (StartTodoUseCase, error) {
	todoRepository, err := do.Invoke[todo.TodoRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke todo repository: %w", err)
	}
	transactionManager, err := do.Invoke[uow.TransactionRunner](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke transaction manager: %w", err)
	}

	return &startTodoUseCase{
		todoRepository: todoRepository,
		txRunner:       transactionManager,
	}, nil
}

// Execute starts a Todo by changing its status to in progress.
func (u *startTodoUseCase) Execute(ctx context.Context, req StartTodoRequest) error {
	err := u.txRunner.RunInTx(ctx, func(ctx context.Context) error {
		foundTodo, err := u.todoRepository.FindByID(ctx, req.ID)
		if err != nil {
			if errors.Is(err, todo.ErrNotFound) {
				return err
			}
			return fmt.Errorf("failed to find todo: %w", err)
		}

		if err := foundTodo.Start(); err != nil {
			// Preserve domain errors (ErrAlreadyCompleted)
			if errors.Is(err, todo.ErrAlreadyCompleted) {
				return err
			}
			return fmt.Errorf("failed to start todo: %w", err)
		}

		if err := u.todoRepository.Update(ctx, foundTodo); err != nil {
			return fmt.Errorf("failed to update todo: %w", err)
		}
		return nil
	})
	if err != nil {
		// Preserve domain errors
		if errors.Is(err, todo.ErrNotFound) || errors.Is(err, todo.ErrAlreadyCompleted) {
			return err
		}
		return fmt.Errorf("failed to execute transaction: %w", err)
	}
	return nil
}
