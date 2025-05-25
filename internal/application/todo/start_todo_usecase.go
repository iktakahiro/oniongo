package todoapp

import (
	"context"
	"fmt"

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
}

// NewStartTodoUseCase creates a new StartTodoUseCase.
func NewStartTodoUseCase(todoRepository todo.TodoRepository) StartTodoUseCase {
	return &startTodoUseCase{todoRepository: todoRepository}
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
	err = u.todoRepository.Update(ctx, todo)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}
