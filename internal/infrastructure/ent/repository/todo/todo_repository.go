package todopsql

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/domain/todo"
	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/db"
	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/entgen"
	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/entgen/todoschema"
)

// todoPsqlRepository is the implementation of the TodoRepository interface.
type todoRepository struct{}

// NewTodoSqliteRepository creates a new TodoSqliteRepository.
func NewTodoSqliteRepository() todo.TodoRepository {
	return &todoRepository{}
}

// Create creates the Todo.
func (r todoRepository) Create(ctx context.Context, todo *todo.Todo) error {
	tx, err := db.GetTx(ctx)
	if err != nil {
		return err
	}

	status := todoschema.Status(todo.Status().String())
	_, err = tx.TodoSchema.Create().
		SetTitle(todo.Title()).
		SetBody(todo.Body()).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}

// FindAll returns all Todos.
func (r todoRepository) FindAll(ctx context.Context) (todos []*todo.Todo, err error) {
	tx, err := db.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	entities, err := tx.TodoSchema.
		Query().
		Where(todoschema.DeletedAtIsNil()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find all todos: %w", err)
	}

	todos = make([]*todo.Todo, len(entities))
	for i, entity := range entities {
		todos[i], err = convertEntToTodo(entity)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to domain %v: %w", entity.ID, err)
		}
	}
	return todos, nil
}

// FindByID returns the Todo with the given ID.
func (r todoRepository) FindByID(
	ctx context.Context,
	id todo.TodoID,
) (todo *todo.Todo, err error) {
	tx, err := db.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	entity, err := tx.TodoSchema.Get(ctx, id.UUID())
	if err != nil {
		return nil, fmt.Errorf("failed to find todo %v: %w", id, err)
	}

	return convertEntToTodo(entity)
}

// Update updates the Todo with the given ID.
func (r todoRepository) Update(ctx context.Context, todo *todo.Todo) (err error) {
	tx, err := db.GetTx(ctx)
	if err != nil {
		return err
	}

	_, err = tx.TodoSchema.UpdateOneID(todo.ID().UUID()).
		SetTitle(todo.Title()).
		SetBody(todo.Body()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update todo %v: %w", todo.ID(), err)
	}
	return nil
}

// Delete deletes the Todo with the given ID.
func (r todoRepository) Delete(ctx context.Context, id todo.TodoID) (err error) {
	tx, err := db.GetTx(ctx)
	if err != nil {
		return err
	}

	err = tx.TodoSchema.DeleteOneID(id.UUID()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete todo %v: %w", id, err)
	}
	return nil
}

// convertEntToTodo converts ent.TodoSchema to domain Todo
func convertEntToTodo(v *entgen.TodoSchema) (*todo.Todo, error) {
	status, err := todo.NewTodoStatusFromString(string(v.Status))
	if err != nil {
		return nil, fmt.Errorf("failed to convert status %v: %w", v.Status, err)
	}
	return todo.ReconstructTodo(v.ID, v.Title, *v.Body, status, v.CreatedAt, v.UpdatedAt), nil
}
