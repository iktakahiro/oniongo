package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/domain/todo"
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
}

// NewCompleteTodoUseCase creates a new CompleteTodoUseCase.
func NewCompleteTodoUseCase(todoRepository todo.TodoRepository) CompleteTodoUseCase {
	return &completeTodoUseCase{todoRepository: todoRepository}
}

// Execute completes a Todo by changing its status to completed.
func (u *completeTodoUseCase) Execute(ctx context.Context, req CompleteTodoRequest) error {
	todo, err := u.todoRepository.FindByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to find todo: %w", err)
	}

	// Check if the todo is already completed
	if todo.IsCompleted() {
		return fmt.Errorf("todo is already completed")
	}

	if err := todo.Complete(); err != nil {
		return fmt.Errorf("failed to complete todo: %w", err)
	}
	err = u.todoRepository.Update(ctx, todo)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}
