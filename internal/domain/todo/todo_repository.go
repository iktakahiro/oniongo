package todo

import (
	"context"
)

// TodoRepository is the interface that wraps the basic CRUD operations for Todo.
type TodoRepository interface {
	Save(ctx context.Context, todo *Todo) error
	FindAll(ctx context.Context) ([]*Todo, error)
	FindByID(ctx context.Context, id TodoID) (*Todo, error)
	Delete(ctx context.Context, id TodoID) error
}
