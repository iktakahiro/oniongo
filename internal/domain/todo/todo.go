// Package todo provides the domain layer for the todo app.
package todo

import (
	"time"

	"github.com/google/uuid"
)

// Todo is the entity that represents a todo item.
type Todo struct {
	id        TodoID
	title     string
	body      string
	createdAt time.Time
	updatedAt time.Time
}

// NewTodo creates a new Todo.
func NewTodo(title string, body string) (*Todo, error) {
	now := time.Now()
	return &Todo{
		id:        NewTodoID(),
		title:     title,
		body:      body,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (t Todo) ID() TodoID {
	return t.id
}

func (t Todo) Title() string {
	return t.title
}

func (t Todo) Body() string {
	return t.body
}

func (t Todo) CreatedAt() time.Time {
	return t.createdAt
}

func (t Todo) UpdatedAt() time.Time {
	return t.updatedAt
}

// ReconstructTodo reconstructs a Todo from the given values.
func ReconstructTodo(id uuid.UUID, title string, body string, createdAt time.Time, updatedAt time.Time) *Todo {
	return &Todo{
		id:        TodoID(id),
		title:     title,
		body:      body,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}