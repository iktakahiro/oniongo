package todo

import "errors"

// Domain errors
var (
	// ErrNotFound is returned when a todo is not found
	ErrNotFound = errors.New("todo not found")
	
	// ErrTitleRequired is returned when title is empty
	ErrTitleRequired = errors.New("title is required")
	
	// ErrInvalidStateTransition is returned when trying to change to an invalid state
	ErrInvalidStateTransition = errors.New("invalid state transition")
	
	// ErrAlreadyCompleted is returned when trying to modify a completed todo
	ErrAlreadyCompleted = errors.New("todo is already completed")
)