package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/domain/todo"
)

type GetTodosRequest struct {}

// GetTodosUseCase is the interface that wraps the basic GetTodos operation.
type GetTodosUseCase interface {
	Execute(ctx context.Context, req GetTodosRequest) ([]*todo.Todo, error)
}

// getTodosUseCase is the implementation of the GetTodosUseCase interface.
type getTodosUseCase struct {
	todoRepository todo.TodoRepository
}

// NewGetTodosUseCase creates a new GetTodosUseCase.
func NewGetTodosUseCase(todoRepository todo.TodoRepository) GetTodosUseCase {
	return &getTodosUseCase{todoRepository: todoRepository}
}

// Execute finds all todos.
func (u getTodosUseCase) Execute(ctx context.Context, req GetTodosRequest) ([]*todo.Todo, error) {
	todos, err := u.todoRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find todos: %w", err)
	}
	return todos, nil
}