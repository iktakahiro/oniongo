// Package todoapp provides the application layer for the todo app.
package todoapp

import (
	"context"
	"fmt"

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
}

// NewCreateTodoUseCase creates a new CreateTodoUseCase.
func NewCreateTodoUseCase(todoRepository todo.TodoRepository) CreateTodoUseCase {
	return &createTodoUseCase{todoRepository: todoRepository}
}

// Execute creates a new Todo.
func (u createTodoUseCase) Execute(ctx context.Context, req CreateTodoRequest) error {
	todo, err := todo.NewTodo(req.Title, req.Body)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	err = u.todoRepository.Save(ctx, todo)
	if err != nil {
		return fmt.Errorf("failed to save todo: %w", err)
	}
	return nil
}
