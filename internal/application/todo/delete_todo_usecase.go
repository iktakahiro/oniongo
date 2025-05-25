package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/domain/todo"
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
}

// NewDeleteTodoUseCase creates a new DeleteTodoUseCase.
func NewDeleteTodoUseCase(todoRepository todo.TodoRepository) DeleteTodoUseCase {
	return &deleteTodoUseCase{todoRepository: todoRepository}
}

// Execute deletes a Todo by its ID.
func (u deleteTodoUseCase) Execute(ctx context.Context, req DeleteTodoRequest) error {
	err := u.todoRepository.Delete(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}
