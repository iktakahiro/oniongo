// Package todo provides the domain layer for the todo app.
package todo

import "time"

// Todo is the entity that represents a todo item.
type Todo struct {
	ID        TodoID
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewTodo creates a new Todo.
func NewTodo(title string, body string) (*Todo, error) {
	now := time.Now()
	return &Todo{
		ID:        NewTodoID(),
		Title:     title,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
