package todo

import "fmt"

// NotFoundError represents an error when a todo is not found
type NotFoundError struct {
	ID TodoID
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("todo not found: %s", e.ID.String())
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// StateError represents an invalid state transition error
type StateError struct {
	Current TodoStatus
	Message string
}

func (e *StateError) Error() string {
	return e.Message
}
