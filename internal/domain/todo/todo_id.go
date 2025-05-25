package todo

import (
	"fmt"

	"github.com/google/uuid"
)

// TodoID is the identifier for a Todo.
type TodoID uuid.UUID

// NewTodoID creates a new TodoID.
func NewTodoID() TodoID {
	id, _ := uuid.NewV7()
	return TodoID(id)
}

// String returns the string representation of the TodoID.
func (id TodoID) String() string {
	return id.UUID().String()
}

// UUID returns the UUID representation of the TodoID.
func (id TodoID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

// NewTodoIDFromString creates a new TodoID from a string.
func NewTodoIDFromString(s string) (TodoID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return TodoID{}, fmt.Errorf("failed to parse uuid %s: %w", s, err)
	}
	return TodoID(id), nil
}
