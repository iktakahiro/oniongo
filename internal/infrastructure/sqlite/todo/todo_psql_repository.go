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

// Save creates the Todo.
func (r todoPsqlRepository) Save(ctx context.Context, todo *todo.Todo) error {
	_, err := r.client.TodoSchema.Create().SetTitle(todo.Title).SetBody(todo.Body).Save(ctx)
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
		todos[i], err = r.toDomain(entity)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to domain %v: %w", entity.ID, err)
		}
	}
	return todos, nil
}

// FindByID returns the Todo with the given ID.
func (r todoPsqlRepository) FindByID(ctx context.Context, id todo.TodoID) (todo *todo.Todo, err error) {
	entity, err := r.client.TodoSchema.Query().Where(
		todoschema.IDEQ(id.UUID()),
	).First(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find todo %v: %w", id, err)
	}

	todo, err = r.toDomain(entity)

	return
}

// Delete deletes the Todo with the given ID.
func (r todoPsqlRepository) Delete(ctx context.Context, id todo.TodoID) (err error) {
	err = r.client.TodoSchema.DeleteOneID(id.UUID()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete todo %v: %w", id, err)
	}
	return nil
}

func (r todoPsqlRepository) toDomain(v *ent.TodoSchema) (*todo.Todo, error) {
	id := todo.TodoID(v.ID)
	return &todo.Todo{
		ID:        id,
		Title:     v.Title,
		Body:      v.Body,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}, nil
}
