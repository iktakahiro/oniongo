package todohandler

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
	pb "github.com/iktakahiro/oniongo/internal/infrastructure/grpc/gen/oniongo/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainTodoToProto(t *testing.T) {
	tests := []struct {
		name     string
		input    *todo.Todo
		expected *pb.Todo
	}{
		{
			name: "converts todo with all fields including completed_at",
			input: func() *todo.Todo {
				createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)
				completedAt := time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC)
				id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

				return todo.ReconstructTodoWithStatus(
					id,
					"Test Todo",
					"Test Body",
					todo.TodoStatusCompleted,
					createdAt,
					updatedAt,
					&completedAt,
				)
			}(),
			expected: &pb.Todo{
				Id:          "550e8400-e29b-41d4-a716-446655440000",
				Title:       "Test Todo",
				Body:        "Test Body",
				Status:      pb.TodoStatus_TODO_STATUS_COMPLETED,
				CreatedAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
				UpdatedAt:   time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC).Unix(),
				CompletedAt: func() *int64 { t := time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC).Unix(); return &t }(),
			},
		},
		{
			name: "converts todo without completed_at",
			input: func() *todo.Todo {
				createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)
				id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

				return todo.ReconstructTodo(
					id,
					"In Progress Todo",
					"In Progress Body",
					todo.TodoStatusInProgress,
					createdAt,
					updatedAt,
				)
			}(),
			expected: &pb.Todo{
				Id:          "550e8400-e29b-41d4-a716-446655440001",
				Title:       "In Progress Todo",
				Body:        "In Progress Body",
				Status:      pb.TodoStatus_TODO_STATUS_IN_PROGRESS,
				CreatedAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
				UpdatedAt:   time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC).Unix(),
				CompletedAt: nil,
			},
		},
		{
			name: "converts todo with not started status",
			input: func() *todo.Todo {
				createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
				id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

				return todo.ReconstructTodo(
					id,
					"Not Started Todo",
					"",
					todo.TodoStatusNotStarted,
					createdAt,
					updatedAt,
				)
			}(),
			expected: &pb.Todo{
				Id:          "550e8400-e29b-41d4-a716-446655440002",
				Title:       "Not Started Todo",
				Body:        "",
				Status:      pb.TodoStatus_TODO_STATUS_NOT_STARTED,
				CreatedAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
				UpdatedAt:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
				CompletedAt: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domainTodoToProto(tt.input)

			assert.Equal(t, tt.expected.Id, result.Id)
			assert.Equal(t, tt.expected.Title, result.Title)
			assert.Equal(t, tt.expected.Body, result.Body)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)

			if tt.expected.CompletedAt == nil {
				assert.Nil(t, result.CompletedAt)
			} else {
				require.NotNil(t, result.CompletedAt)
				assert.Equal(t, *tt.expected.CompletedAt, *result.CompletedAt)
			}
		})
	}
}

func TestDomainStatusToProtoStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    todo.TodoStatus
		expected pb.TodoStatus
	}{
		{
			name:     "converts not started status",
			input:    todo.TodoStatusNotStarted,
			expected: pb.TodoStatus_TODO_STATUS_NOT_STARTED,
		},
		{
			name:     "converts in progress status",
			input:    todo.TodoStatusInProgress,
			expected: pb.TodoStatus_TODO_STATUS_IN_PROGRESS,
		},
		{
			name:     "converts completed status",
			input:    todo.TodoStatusCompleted,
			expected: pb.TodoStatus_TODO_STATUS_COMPLETED,
		},
		{
			name:     "converts invalid status to unspecified",
			input:    todo.TodoStatus(999), // invalid status
			expected: pb.TodoStatus_TODO_STATUS_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domainStatusToProtoStatus(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseUUIDFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    todo.TodoID
	}{
		{
			name:        "parses valid UUID",
			input:       "550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
			expected:    todo.TodoID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
		},
		{
			name:        "parses valid UUID with different format",
			input:       "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			expectError: false,
			expected:    todo.TodoID(uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")),
		},
		{
			name:        "returns error for invalid UUID format",
			input:       "invalid-uuid",
			expectError: true,
			expected:    todo.TodoID{},
		},
		{
			name:        "returns error for empty string",
			input:       "",
			expectError: true,
			expected:    todo.TodoID{},
		},
		{
			name:        "parses UUID without hyphens",
			input:       "550e8400e29b41d4a716446655440000",
			expectError: false,
			expected:    todo.TodoID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
		},
		{
			name:        "returns error for partial UUID",
			input:       "550e8400-e29b-41d4",
			expectError: true,
			expected:    todo.TodoID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseUUIDFromString(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, todo.TodoID{}, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Benchmark tests for performance validation
func BenchmarkDomainTodoToProto(b *testing.B) {
	createdAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)
	completedAt := time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC)
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	domainTodo := todo.ReconstructTodoWithStatus(
		id,
		"Benchmark Todo",
		"Benchmark Body",
		todo.TodoStatusCompleted,
		createdAt,
		updatedAt,
		&completedAt,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = domainTodoToProto(domainTodo)
	}
}

func BenchmarkDomainStatusToProtoStatus(b *testing.B) {
	status := todo.TodoStatusInProgress

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = domainStatusToProtoStatus(status)
	}
}

func BenchmarkParseUUIDFromString(b *testing.B) {
	uuidStr := "550e8400-e29b-41d4-a716-446655440000"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parseUUIDFromString(uuidStr)
	}
}
