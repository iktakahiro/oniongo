package todopsql

import (
	"context"
	"fmt"

	"github.com/iktakahiro/oniongo/internal/domain/todo"
	"github.com/iktakahiro/oniongo/internal/infrastructure/sqlite/ent"
	"github.com/iktakahiro/oniongo/internal/infrastructure/sqlite/ent/todoschema"
)

// todoPsqlRepository is the implementation of the TodoRepository interface.
type todoPsqlRepository struct {
	client *ent.Client
}

// NewTodoPsqlRepository creates a new TodoPsqlRepository.
func NewTodoPsqlRepository() todo.TodoRepository {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil
	}

	return &todoPsqlRepository{client: client}
}

// Create creates the Todo.
func (r todoPsqlRepository) Create(ctx context.Context, todo *todo.Todo) error {
	status := todoschema.Status(todo.Status().String())
	_, err := r.client.TodoSchema.Create().
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
func (r todoPsqlRepository) FindAll(ctx context.Context) (todos []*todo.Todo, err error) {
	entities, err := r.client.TodoSchema.Query().All(ctx)
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
func (r todoPsqlRepository) FindByID(ctx context.Context, id todo.TodoID) (todo *todo.Todo, err error) {
	entity, err := r.client.TodoSchema.Get(ctx, id.UUID())
	if err != nil {
		return nil, fmt.Errorf("failed to find todo %v: %w", id, err)
	}

	return convertEntToTodo(entity)
}

// Update updates the Todo with the given ID.
func (r todoPsqlRepository) Update(ctx context.Context, todo *todo.Todo) (err error) {
	_, err = r.client.TodoSchema.UpdateOneID(todo.ID().UUID()).SetTitle(todo.Title()).SetBody(todo.Body()).Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update todo %v: %w", todo.ID(), err)
	}
	return nil
}
// Delete deletes the Todo with the given ID.
func (r todoPsqlRepository) Delete(ctx context.Context, id todo.TodoID) (err error) {
	err = r.client.TodoSchema.DeleteOneID(id.UUID()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete todo %v: %w", id, err)
	}
	return nil
}

// convertEntToTodo converts ent.TodoSchema to domain Todo
func convertEntToTodo(v *ent.TodoSchema) (*todo.Todo, error) {
	status, err := todo.NewTodoStatusFromString(string(v.Status))
	if err != nil {
		return nil, fmt.Errorf("failed to convert status %v: %w", v.Status, err)
	}
	return todo.ReconstructTodo(v.ID, v.Title, v.Body, status, v.CreatedAt, v.UpdatedAt), nil
}
