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

// ID returns the ID of the Todo.
func (t Todo) ID() TodoID {
	return t.id
}

// Title returns the title of the Todo.
func (t Todo) Title() string {
	return t.title
}

// Body returns the body of the Todo.
func (t Todo) Body() string {
	return t.body
}

// CreatedAt returns the created at of the Todo.
func (t Todo) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the updated at of the Todo.
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