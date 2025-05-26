package todoapp

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/application"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
)

type UpdateTodoRequest struct {
	ID    todo.TodoID
	Title string
	Body  string
}

// UpdateTodoUseCase is the interface that wraps the basic UpdateTodo operation.
type UpdateTodoUseCase interface {
	Execute(ctx context.Context, req UpdateTodoRequest) error
}

// updateTodoUseCase is the implementation of the UpdateTodoUseCase interface.
type updateTodoUseCase struct {
	todoRepository todo.TodoRepository
	txManager      application.TransactionManager
}

// NewUpdateTodoUseCase creates a new UpdateTodoUseCase.
func NewUpdateTodoUseCase(
	todoRepository todo.TodoRepository,
	transactionManager application.TransactionManager,
) UpdateTodoUseCase {
	return &updateTodoUseCase{
		todoRepository: todoRepository,
		txManager:      transactionManager,
	}
}

// Execute updates a Todo by its ID.
func (u *updateTodoUseCase) Execute(ctx context.Context, req UpdateTodoRequest) error {
	todo, err := u.todoRepository.FindByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("failed to find todo: %w", err)
	}
	if err := todo.SetTitle(req.Title); err != nil {
		return fmt.Errorf("failed to set title: %w", err)
	}
	if err := todo.SetBody(req.Body); err != nil {
		return fmt.Errorf("failed to set body: %w", err)
	}
	err = u.txManager.RunInTx(ctx, func(ctx context.Context) error {
		if err := u.todoRepository.Update(ctx, todo); err != nil {
			return fmt.Errorf("failed to update todo: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}
	return nil
}
