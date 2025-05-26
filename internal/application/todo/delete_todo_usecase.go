package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
	"github.com/samber/do"
)

type DeleteTodoRequest struct {
	ID todo.TodoID
}

// DeleteTodoUseCase is the interface that wraps the basic DeleteTodo operation.
type DeleteTodoUseCase interface {
	Execute(ctx context.Context, req DeleteTodoRequest) error
}

// deleteTodoUseCase is the implementation of the DeleteTodoUseCase interface.
type deleteTodoUseCase struct {
	todoRepository todo.TodoRepository
	txManager      application.TransactionManager
}

// NewDeleteTodoUseCase creates a new DeleteTodoUseCase.
func NewDeleteTodoUseCase(i *do.Injector) (DeleteTodoUseCase, error) {
	todoRepository, err := do.Invoke[todo.TodoRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke todo repository: %w", err)
	}
	transactionManager, err := do.Invoke[application.TransactionManager](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke transaction manager: %w", err)
	}

	return &deleteTodoUseCase{
		todoRepository: todoRepository,
		txManager:      transactionManager,
	}, nil
}

// Execute deletes a Todo by its ID.
func (u deleteTodoUseCase) Execute(ctx context.Context, req DeleteTodoRequest) error {
	err := u.txManager.RunInTx(ctx, func(ctx context.Context) error {
		if err := u.todoRepository.Delete(ctx, req.ID); err != nil {
			return fmt.Errorf("failed to delete todo: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}
	return nil
}
