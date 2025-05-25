// Package todo provides the domain layer for the todo app.
package todo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Todo is the entity that represents a todo item.
type Todo struct {
	id          TodoID
	title       string
	body        string
	status      TodoStatus
	createdAt   time.Time
	updatedAt   time.Time
	completedAt *time.Time
}

// NewTodo creates a new Todo.
func NewTodo(title string, body string) (*Todo, error) {
	now := time.Now()
	return &Todo{
		id:          NewTodoID(),
		title:       title,
		body:        body,
		status:      TodoStatusNotStarted,
		createdAt:   now,
		updatedAt:   now,
		completedAt: nil,
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

// Status returns the status of the Todo.
func (t Todo) Status() TodoStatus {
	return t.status
}

// CreatedAt returns the created at of the Todo.
func (t Todo) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the updated at of the Todo.
func (t Todo) UpdatedAt() time.Time {
	return t.updatedAt
}

// CompletedAt returns the completed at of the Todo.
func (t Todo) CompletedAt() *time.Time {
	return t.completedAt
}

func (t *Todo) SetTitle(title string) error {
	if title == "" {
		return errors.New("title is required")
	}
	t.title = title
	t.updatedAt = time.Now()
	return nil
}

func (t *Todo) SetBody(body string) error {
	t.body = body
	t.updatedAt = time.Now()
	return nil
}

// Start changes the Todo's status to in progress.
func (t *Todo) Start() error {
	if t.status == TodoStatusCompleted {
		return errors.New("cannot start a completed todo")
	}
	t.status = TodoStatusInProgress
	t.updatedAt = time.Now()
	return nil
}

// Complete changes the Todo's status to completed.
func (t *Todo) Complete() error {
	if t.status == TodoStatusCompleted {
		return errors.New("todo is already completed")
	}
	now := time.Now()
	t.status = TodoStatusCompleted
	t.completedAt = &now
	t.updatedAt = now
	return nil
}

// IsInProgress checks if the Todo is in progress.
func (t Todo) IsInProgress() bool {
	return t.status == TodoStatusInProgress
}

// IsCompleted checks if the Todo is completed.
func (t Todo) IsCompleted() bool {
	return t.status == TodoStatusCompleted
}

// ReconstructTodo reconstructs a Todo from the given values.
func ReconstructTodo(id uuid.UUID, title string, body string, createdAt time.Time, updatedAt time.Time) *Todo {
	return &Todo{
		id:          TodoID(id),
		title:       title,
		body:        body,
		status:      TodoStatusNotStarted,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		completedAt: nil,
	}
}

// ReconstructTodoWithStatus reconstructs a Todo from the given values including status and completedAt.
func ReconstructTodoWithStatus(id uuid.UUID, title string, body string, status TodoStatus, createdAt time.Time, updatedAt time.Time, completedAt *time.Time) *Todo {
	return &Todo{
		id:          TodoID(id),
		title:       title,
		body:        body,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		completedAt: completedAt,
	}
}