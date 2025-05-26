package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
)

type GetTodosRequest struct{}

// GetTodosUseCase is the interface that wraps the basic GetTodos operation.
type GetTodosUseCase interface {
	Execute(ctx context.Context, req GetTodosRequest) ([]*todo.Todo, error)
}

// getTodosUseCase is the implementation of the GetTodosUseCase interface.
type getTodosUseCase struct {
	todoRepository todo.TodoRepository
	txManager      application.TransactionManager
}

// NewGetTodosUseCase creates a new GetTodosUseCase.
func NewGetTodosUseCase(
	todoRepository todo.TodoRepository,
	transactionManager application.TransactionManager,
) GetTodosUseCase {
	return &getTodosUseCase{
		todoRepository: todoRepository,
		txManager:      transactionManager,
	}
}

// Execute finds all todos.
func (u getTodosUseCase) Execute(ctx context.Context, req GetTodosRequest) ([]*todo.Todo, error) {
	var result []*todo.Todo
	err := u.txManager.RunInTx(ctx, func(ctx context.Context) error {
		todos, err := u.todoRepository.FindAll(ctx)
		if err != nil {
			return fmt.Errorf("failed to find todos: %w", err)
		}
		result = todos
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}
	return result, nil
}
