package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
)

type GetTodoRequest struct {
	ID todo.TodoID
}

// GetTodoUseCase is the interface that wraps the basic GetTodo operation.
type GetTodoUseCase interface {
	Execute(ctx context.Context, req GetTodoRequest) (*todo.Todo, error)
}

// getTodoUseCase is the implementation of the GetTodoUseCase interface.
type getTodoUseCase struct {
	todoRepository todo.TodoRepository
	txManager      application.TransactionManager
}

// NewGetTodoUseCase creates a new GetTodoUseCase.
func NewGetTodoUseCase(
	todoRepository todo.TodoRepository,
	transactionManager application.TransactionManager,
) GetTodoUseCase {
	return &getTodoUseCase{
		todoRepository: todoRepository,
		txManager:      transactionManager,
	}
}

// Execute gets a Todo by its ID.
func (u getTodoUseCase) Execute(ctx context.Context, req GetTodoRequest) (*todo.Todo, error) {
	var result *todo.Todo
	err := u.txManager.RunInTx(ctx, func(ctx context.Context) error {
		todo, err := u.todoRepository.FindByID(ctx, req.ID)
		if err != nil {
			return fmt.Errorf("failed to find todo: %w", err)
		}
		result = todo
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}
	return result, nil
}
