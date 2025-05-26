package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application/uow"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
	"github.com/samber/do"
)

type GetTodosRequest struct{}

// GetTodosUseCase is the interface that wraps the basic GetTodos operation.
type GetTodosUseCase interface {
	Execute(ctx context.Context, req GetTodosRequest) ([]*todo.Todo, error)
}

// getTodosUseCase is the implementation of the GetTodosUseCase interface.
type getTodosUseCase struct {
	todoRepository todo.TodoRepository
	txRunner       uow.TransactionRunner
}

// NewGetTodosUseCase creates a new GetTodosUseCase.
func NewGetTodosUseCase(i *do.Injector) (GetTodosUseCase, error) {
	todoRepository, err := do.Invoke[todo.TodoRepository](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke todo repository: %w", err)
	}
	transactionManager, err := do.Invoke[uow.TransactionRunner](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke transaction manager: %w", err)
	}

	return &getTodosUseCase{
		todoRepository: todoRepository,
		txRunner:       transactionManager,
	}, nil
}

// Execute finds all todos.
func (u getTodosUseCase) Execute(ctx context.Context, req GetTodosRequest) ([]*todo.Todo, error) {
	var result []*todo.Todo
	err := u.txRunner.RunInTx(ctx, func(ctx context.Context) error {
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
