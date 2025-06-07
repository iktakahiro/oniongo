package todo

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		body        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid todo with title and body",
			title:       "Test Todo",
			body:        "This is a test todo",
			expectError: false,
		},
		{
			name:        "todo with empty body",
			title:       "Test Todo",
			body:        "",
			expectError: false,
		},
		{
			name:        "todo with empty title",
			title:       "",
			body:        "This is a test todo",
			expectError: true,
			errorMsg:    "title: title is required",
		},
		{
			name:        "todo with both empty",
			title:       "",
			body:        "",
			expectError: true,
			errorMsg:    "title: title is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			before := time.Now()

			// When
			todo, err := NewTodo(tt.title, tt.body)

			// Then
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.errorMsg, err.Error())
				require.Nil(t, todo)
			} else {
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
			}
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

func TestReconstructTodoWithStatus(t *testing.T) {
	tests := []struct {
		name        string
		id          uuid.UUID
		title       string
		body        string
		status      TodoStatus
		createdAt   time.Time
		updatedAt   time.Time
		completedAt *time.Time
	}{
		{
			name:        "reconstruction with completed status",
			id:          uuid.New(),
			title:       "Completed Todo",
			body:        "This is a completed todo",
			status:      TodoStatusCompleted,
			createdAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			updatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			completedAt: func() *time.Time { t := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC); return &t }(),
		},
		{
			name:        "reconstruction with not started status",
			id:          uuid.New(),
			title:       "Not Started Todo",
			body:        "This is a not started todo",
			status:      TodoStatusNotStarted,
			createdAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			updatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			completedAt: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			todo := ReconstructTodoWithStatus(
				tt.id,
				tt.title,
				tt.body,
				tt.status,
				tt.createdAt,
				tt.updatedAt,
				tt.completedAt,
			)

			// Then
			require.NotNil(t, todo)
			require.Equal(t, TodoID(tt.id), todo.ID())
			require.Equal(t, tt.title, todo.Title())
			require.Equal(t, tt.body, todo.Body())
			require.Equal(t, tt.status, todo.Status())
			require.Equal(t, tt.createdAt, todo.CreatedAt())
			require.Equal(t, tt.updatedAt, todo.UpdatedAt())
			require.Equal(t, tt.completedAt, todo.CompletedAt())
		})
	}
}

func TestTodo_SetTitle(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid title",
			title:       "New Title",
			expectError: false,
		},
		{
			name:        "empty title",
			title:       "",
			expectError: true,
			errorMsg:    "title: title is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			todo, err := NewTodo("Original Title", "Original Body")
			require.NoError(t, err)
			originalUpdatedAt := todo.UpdatedAt()
			time.Sleep(1 * time.Millisecond) // Ensure time difference

			// When
			err = todo.SetTitle(tt.title)

			// Then
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.errorMsg, err.Error())
				require.Equal(t, "Original Title", todo.Title())      // Title should not change
				require.Equal(t, originalUpdatedAt, todo.UpdatedAt()) // UpdatedAt should not change
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.title, todo.Title())
				require.True(t, todo.UpdatedAt().After(originalUpdatedAt))
			}
		})
	}
}

func TestTodo_SetBody(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "valid body",
			body: "New Body",
		},
		{
			name: "empty body",
			body: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			todo, err := NewTodo("Original Title", "Original Body")
			require.NoError(t, err)
			originalUpdatedAt := todo.UpdatedAt()
			time.Sleep(1 * time.Millisecond) // Ensure time difference

			// When
			err = todo.SetBody(tt.body)

			// Then
			require.NoError(t, err)
			require.Equal(t, tt.body, todo.Body())
			require.True(t, todo.UpdatedAt().After(originalUpdatedAt))
		})
	}
}

func TestTodo_Start(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus TodoStatus
		expectError   bool
		errorMsg      string
	}{
		{
			name:          "start not started todo",
			initialStatus: TodoStatusNotStarted,
			expectError:   false,
		},
		{
			name:          "start in progress todo",
			initialStatus: TodoStatusInProgress,
			expectError:   false,
		},
		{
			name:          "start completed todo",
			initialStatus: TodoStatusCompleted,
			expectError:   true,
			errorMsg:      "todo is already completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			todo, err := NewTodo("Test Todo", "Test Body")
			require.NoError(t, err)

			// Set initial status
			if tt.initialStatus == TodoStatusCompleted {
				err = todo.Complete()
				require.NoError(t, err)
			} else if tt.initialStatus == TodoStatusInProgress {
				err = todo.Start()
				require.NoError(t, err)
			}

			originalUpdatedAt := todo.UpdatedAt()
			time.Sleep(1 * time.Millisecond) // Ensure time difference

			// When
			err = todo.Start()

			// Then
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.errorMsg, err.Error())
				require.Equal(t, tt.initialStatus, todo.Status())
			} else {
				require.NoError(t, err)
				require.Equal(t, TodoStatusInProgress, todo.Status())
				require.True(t, todo.UpdatedAt().After(originalUpdatedAt))
			}
		})
	}
}

func TestTodo_Complete(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus TodoStatus
		expectError   bool
		errorMsg      string
	}{
		{
			name:          "complete not started todo",
			initialStatus: TodoStatusNotStarted,
			expectError:   false,
		},
		{
			name:          "complete in progress todo",
			initialStatus: TodoStatusInProgress,
			expectError:   false,
		},
		{
			name:          "complete already completed todo",
			initialStatus: TodoStatusCompleted,
			expectError:   true,
			errorMsg:      "todo is already completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			todo, err := NewTodo("Test Todo", "Test Body")
			require.NoError(t, err)

			// Set initial status
			if tt.initialStatus == TodoStatusInProgress {
				err = todo.Start()
				require.NoError(t, err)
			} else if tt.initialStatus == TodoStatusCompleted {
				err = todo.Complete()
				require.NoError(t, err)
			}

			originalUpdatedAt := todo.UpdatedAt()
			originalCompletedAt := todo.CompletedAt()
			time.Sleep(1 * time.Millisecond) // Ensure time difference

			// When
			err = todo.Complete()

			// Then
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.errorMsg, err.Error())
				require.Equal(t, tt.initialStatus, todo.Status())
				require.Equal(t, originalCompletedAt, todo.CompletedAt())
			} else {
				require.NoError(t, err)
				require.Equal(t, TodoStatusCompleted, todo.Status())
				require.True(t, todo.UpdatedAt().After(originalUpdatedAt))
				require.NotNil(t, todo.CompletedAt())
				require.True(t, todo.CompletedAt().After(originalUpdatedAt))
			}
		})
	}
}

func TestTodo_IsInProgress(t *testing.T) {
	tests := []struct {
		name     string
		status   TodoStatus
		expected bool
	}{
		{
			name:     "not started todo",
			status:   TodoStatusNotStarted,
			expected: false,
		},
		{
			name:     "in progress todo",
			status:   TodoStatusInProgress,
			expected: true,
		},
		{
			name:     "completed todo",
			status:   TodoStatusCompleted,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			todo, err := NewTodo("Test Todo", "Test Body")
			require.NoError(t, err)

			// Set status
			switch tt.status {
			case TodoStatusInProgress:
				err = todo.Start()
				require.NoError(t, err)
			case TodoStatusCompleted:
				err = todo.Complete()
				require.NoError(t, err)
			}

			// When & Then
			require.Equal(t, tt.expected, todo.IsInProgress())
		})
	}
}

func TestTodo_IsCompleted(t *testing.T) {
	tests := []struct {
		name     string
		status   TodoStatus
		expected bool
	}{
		{
			name:     "not started todo",
			status:   TodoStatusNotStarted,
			expected: false,
		},
		{
			name:     "in progress todo",
			status:   TodoStatusInProgress,
			expected: false,
		},
		{
			name:     "completed todo",
			status:   TodoStatusCompleted,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			todo, err := NewTodo("Test Todo", "Test Body")
			require.NoError(t, err)

			// Set status
			switch tt.status {
			case TodoStatusInProgress:
				err = todo.Start()
				require.NoError(t, err)
			case TodoStatusCompleted:
				err = todo.Complete()
				require.NoError(t, err)
			}

			// When & Then
			require.Equal(t, tt.expected, todo.IsCompleted())
		})
	}
}
