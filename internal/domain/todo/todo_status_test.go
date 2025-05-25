package todo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTodoStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   TodoStatus
		expected string
	}{
		{
			name:     "not started status",
			status:   TodoStatusNotStarted,
			expected: "NOT_STARTED",
		},
		{
			name:     "in progress status",
			status:   TodoStatusInProgress,
			expected: "IN_PROGRESS",
		},
		{
			name:     "completed status",
			status:   TodoStatusCompleted,
			expected: "COMPLETED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			result := tt.status.String()

			// Then
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestTodoStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   TodoStatus
		expected bool
	}{
		{
			name:     "valid not started status",
			status:   TodoStatusNotStarted,
			expected: true,
		},
		{
			name:     "valid in progress status",
			status:   TodoStatusInProgress,
			expected: true,
		},
		{
			name:     "valid completed status",
			status:   TodoStatusCompleted,
			expected: true,
		},
		{
			name:     "invalid status",
			status:   TodoStatus(999), // Use an invalid int value
			expected: false,
		},
		{
			name:     "negative status",
			status:   TodoStatus(-1), // Use a negative int value
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			result := tt.status.IsValid()

			// Then
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestNewTodoStatusFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    TodoStatus
		expectError bool
	}{
		{
			name:        "valid not started status",
			input:       "NOT_STARTED",
			expected:    TodoStatusNotStarted,
			expectError: false,
		},
		{
			name:        "valid in progress status",
			input:       "IN_PROGRESS",
			expected:    TodoStatusInProgress,
			expectError: false,
		},
		{
			name:        "valid completed status",
			input:       "COMPLETED",
			expected:    TodoStatusCompleted,
			expectError: false,
		},
		{
			name:        "invalid status",
			input:       "INVALID",
			expected:    TodoStatusNotStarted, // Function returns TodoStatusNotStarted on error
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    TodoStatusNotStarted, // Function returns TodoStatusNotStarted on error
			expectError: true,
		},
		{
			name:        "lowercase status",
			input:       "not_started",
			expected:    TodoStatusNotStarted, // Function returns TodoStatusNotStarted on error
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			result, err := NewTodoStatusFromString(tt.input)

			// Then
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.expected, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}
