package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application"
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
	txManager      application.TransactionManager
}

// NewCompleteTodoUseCase creates a new CompleteTodoUseCase.
func NewCompleteTodoUseCase(i *do.Injector) (CompleteTodoUseCase, error) {
	todoRepository, err := do.Invoke[todo.TodoRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke todo repository: %w", err)
	}
	transactionManager, err := do.Invoke[application.TransactionManager](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke transaction manager: %w", err)
	}

	return &completeTodoUseCase{
		todoRepository: todoRepository,
		txManager:      transactionManager,
	}, nil
}

// Execute completes a Todo by changing its status to completed.
func (u *completeTodoUseCase) Execute(ctx context.Context, req CompleteTodoRequest) error {
	err := u.txManager.RunInTx(ctx, func(ctx context.Context) error {
		todo, err := u.todoRepository.FindByID(ctx, req.ID)
		if err != nil {
			return fmt.Errorf("failed to find todo: %w", err)
		}

		if todo.IsCompleted() {
			return fmt.Errorf("todo is already completed")
		}

		if err := todo.Complete(); err != nil {
			return fmt.Errorf("failed to complete todo: %w", err)
		}

		if err := u.todoRepository.Update(ctx, todo); err != nil {
			return fmt.Errorf("failed to update todo: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}
	return nil
}
