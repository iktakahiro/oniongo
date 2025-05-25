package todo

import "fmt"

// TodoStatus represents the status of a Todo.
type TodoStatus int

const (
	// TodoStatusNotStarted represents a todo that has not been started.
	TodoStatusNotStarted TodoStatus = iota
	// TodoStatusInProgress represents a todo that is in progress.
	TodoStatusInProgress
	// TodoStatusCompleted represents a todo that has been completed.
	TodoStatusCompleted
)

// statusStrings maps TodoStatus values to their string representations.
var statusStrings = map[TodoStatus]string{
	TodoStatusNotStarted: "NOT_STARTED",
	TodoStatusInProgress: "IN_PROGRESS",
	TodoStatusCompleted:  "COMPLETED",
}

// stringToStatus maps string representations to TodoStatus values.
var stringToStatus = map[string]TodoStatus{
	"NOT_STARTED": TodoStatusNotStarted,
	"IN_PROGRESS": TodoStatusInProgress,
	"COMPLETED":   TodoStatusCompleted,
}

// String returns the string representation of the TodoStatus.
func (s TodoStatus) String() string {
	if str, ok := statusStrings[s]; ok {
		return str
	}
	return fmt.Sprintf("TodoStatus(%d)", int(s))
}

// IsValid checks if the TodoStatus is valid.
func (s TodoStatus) IsValid() bool {
	_, ok := statusStrings[s]
	return ok
}

// NewTodoStatusFromString creates a TodoStatus from a string.
func NewTodoStatusFromString(s string) (TodoStatus, error) {
	if status, ok := stringToStatus[s]; ok {
		return status, nil
	}
	return TodoStatusNotStarted, fmt.Errorf("invalid todo status: %s", s)
}

// AllStatuses returns all valid TodoStatus values.
func AllStatuses() []TodoStatus {
	return []TodoStatus{
		TodoStatusNotStarted,
		TodoStatusInProgress,
		TodoStatusCompleted,
	}
}
