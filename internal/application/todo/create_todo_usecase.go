// Package todoapp provides the application layer for the todo app.
package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
)

type CreateTodoRequest struct {
	Title string
	Body  string
}

// CreateTodoUseCase is the interface that wraps the basic CreateTodo operation.
type CreateTodoUseCase interface {
	Execute(ctx context.Context, req CreateTodoRequest) error
}

// createTodoUseCase is the implementation of the CreateTodoUseCase interface.
type createTodoUseCase struct {
	todoRepository todo.TodoRepository
	txManager      application.TransactionManager
}

// NewCreateTodoUseCase creates a new CreateTodoUseCase.
func NewCreateTodoUseCase(
	todoRepository todo.TodoRepository,
	transactionManager application.TransactionManager,
) CreateTodoUseCase {
	return &createTodoUseCase{
		todoRepository: todoRepository,
		txManager:      transactionManager,
	}
}

// Execute creates a new Todo.
func (u createTodoUseCase) Execute(ctx context.Context, req CreateTodoRequest) error {
	todo, err := todo.NewTodo(req.Title, req.Body)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	err = u.txManager.RunInTx(ctx, func(ctx context.Context) error {
		if err := u.todoRepository.Create(ctx, todo); err != nil {
			return fmt.Errorf("failed to save todo: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}
	return nil
}
