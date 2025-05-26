package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
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
	txManager      application.TransactionManager
}

// NewStartTodoUseCase creates a new StartTodoUseCase.
func NewStartTodoUseCase(
	todoRepository todo.TodoRepository,
	transactionManager application.TransactionManager,
) StartTodoUseCase {
	return &startTodoUseCase{
		todoRepository: todoRepository,
		txManager:      transactionManager,
	}
}

// Execute starts a Todo by changing its status to in progress.
func (u *startTodoUseCase) Execute(ctx context.Context, req StartTodoRequest) error {
	todo, err := u.todoRepository.FindByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to find todo: %w", err)
	}
	if err := todo.Start(); err != nil {
		return fmt.Errorf("failed to start todo: %w", err)
	}
	err = u.txManager.RunInTx(ctx, func(ctx context.Context) error {
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
