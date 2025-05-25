package todo

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name  string
		title string
		body  string
	}{
		{
			name:  "valid todo with title and body",
			title: "Test Todo",
			body:  "This is a test todo",
		},
		{
			name:  "todo with empty body",
			title: "Test Todo",
			body:  "",
		},
		{
			name:  "todo with empty title",
			title: "",
			body:  "This is a test todo",
		},
		{
			name:  "todo with both empty",
			title: "",
			body:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			before := time.Now()

			// When
			todo, err := NewTodo(tt.title, tt.body)

			// Then
			require.NoError(t, err)
			require.NotNil(t, todo)
			require.Equal(t, tt.title, todo.Title())
			require.Equal(t, tt.body, todo.Body())
			require.NotEqual(t, TodoID{}, todo.ID())
			require.False(t, todo.CreatedAt().IsZero())
			require.False(t, todo.UpdatedAt().IsZero())
			require.True(t, todo.CreatedAt().After(before) || todo.CreatedAt().Equal(before))
			require.True(t, todo.UpdatedAt().After(before) || todo.UpdatedAt().Equal(before))
			require.Equal(t, todo.CreatedAt(), todo.UpdatedAt())
		})
	}
}

func TestTodo_Getters(t *testing.T) {
	// Given
	title := "Test Todo"
	body := "This is a test todo"
	todo, err := NewTodo(title, body)
	require.NoError(t, err)

	// When & Then
	require.Equal(t, title, todo.Title())
	require.Equal(t, body, todo.Body())
	require.NotEqual(t, TodoID{}, todo.ID())
	require.False(t, todo.CreatedAt().IsZero())
	require.False(t, todo.UpdatedAt().IsZero())
}

func TestReconstructTodo(t *testing.T) {
	tests := []struct {
		name      string
		id        uuid.UUID
		title     string
		body      string
		status    TodoStatus
		createdAt time.Time
		updatedAt time.Time
	}{
		{
			name:      "valid reconstruction",
			id:        uuid.New(),
			title:     "Reconstructed Todo",
			body:      "This is a reconstructed todo",
			status:    TodoStatusNotStarted,
			createdAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			updatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "reconstruction with empty values",
			id:        uuid.New(),
			title:     "",
			body:      "",
			status:    TodoStatusNotStarted,
			createdAt: time.Time{},
			updatedAt: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			todo := ReconstructTodo(tt.id, tt.title, tt.body, tt.status, tt.createdAt, tt.updatedAt)

			// Then
			require.NotNil(t, todo)
			require.Equal(t, TodoID(tt.id), todo.ID())
			require.Equal(t, tt.title, todo.Title())
			require.Equal(t, tt.body, todo.Body())
			require.Equal(t, tt.status, todo.Status())
			require.Equal(t, tt.createdAt, todo.CreatedAt())
			require.Equal(t, tt.updatedAt, todo.UpdatedAt())
		})
	}
}
