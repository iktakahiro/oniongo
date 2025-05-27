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
		name      string
		setupTodo func() *todo.Todo
		expected  func(*todo.Todo) *pb.Todo
	}{
		{
			name: "converts todo with all fields",
			setupTodo: func() *todo.Todo {
				id := uuid.New()
				createdAt := time.Now().UTC()
				updatedAt := createdAt.Add(time.Hour)
				completedAt := updatedAt.Add(time.Hour)

				todoItem := todo.ReconstructTodoWithStatus(
					id,
					"Test Title",
					"Test Body",
					todo.TodoStatusCompleted,
					createdAt,
					updatedAt,
					&completedAt,
				)
				return todoItem
			},
			expected: func(domainTodo *todo.Todo) *pb.Todo {
				completedAt := domainTodo.CompletedAt().Unix()
				return &pb.Todo{
					Id:          domainTodo.ID().String(),
					Title:       domainTodo.Title(),
					Body:        domainTodo.Body(),
					Status:      pb.TodoStatus_TODO_STATUS_COMPLETED,
					CreatedAt:   domainTodo.CreatedAt().Unix(),
					UpdatedAt:   domainTodo.UpdatedAt().Unix(),
					CompletedAt: &completedAt,
				}
			},
		},
		{
			name: "converts todo without completed_at",
			setupTodo: func() *todo.Todo {
				id := uuid.New()
				createdAt := time.Now().UTC()
				updatedAt := createdAt.Add(time.Hour)

				todoItem := todo.ReconstructTodoWithStatus(
					id,
					"Test Title",
					"Test Body",
					todo.TodoStatusInProgress,
					createdAt,
					updatedAt,
					nil,
				)
				return todoItem
			},
			expected: func(domainTodo *todo.Todo) *pb.Todo {
				return &pb.Todo{
					Id:          domainTodo.ID().String(),
					Title:       domainTodo.Title(),
					Body:        domainTodo.Body(),
					Status:      pb.TodoStatus_TODO_STATUS_IN_PROGRESS,
					CreatedAt:   domainTodo.CreatedAt().Unix(),
					UpdatedAt:   domainTodo.UpdatedAt().Unix(),
					CompletedAt: nil,
				}
			},
		},
		{
			name: "converts todo with not started status",
			setupTodo: func() *todo.Todo {
				id := uuid.New()
				createdAt := time.Now().UTC()
				updatedAt := createdAt

				todoItem := todo.ReconstructTodoWithStatus(
					id,
					"Not Started Todo",
					"Description",
					todo.TodoStatusNotStarted,
					createdAt,
					updatedAt,
					nil,
				)
				return todoItem
			},
			expected: func(domainTodo *todo.Todo) *pb.Todo {
				return &pb.Todo{
					Id:          domainTodo.ID().String(),
					Title:       domainTodo.Title(),
					Body:        domainTodo.Body(),
					Status:      pb.TodoStatus_TODO_STATUS_NOT_STARTED,
					CreatedAt:   domainTodo.CreatedAt().Unix(),
					UpdatedAt:   domainTodo.UpdatedAt().Unix(),
					CompletedAt: nil,
				}
			},
		},
		{
			name: "converts todo created with NewTodo",
			setupTodo: func() *todo.Todo {
				todoItem, err := todo.NewTodo("Simple Todo", "Simple Body")
				require.NoError(t, err)
				return todoItem
			},
			expected: func(domainTodo *todo.Todo) *pb.Todo {
				return &pb.Todo{
					Id:          domainTodo.ID().String(),
					Title:       domainTodo.Title(),
					Body:        domainTodo.Body(),
					Status:      pb.TodoStatus_TODO_STATUS_NOT_STARTED,
					CreatedAt:   domainTodo.CreatedAt().Unix(),
					UpdatedAt:   domainTodo.UpdatedAt().Unix(),
					CompletedAt: nil,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domainTodo := tt.setupTodo()
			expected := tt.expected(domainTodo)

			result := domainTodoToProto(domainTodo)

			assert.Equal(t, expected.Id, result.Id)
			assert.Equal(t, expected.Title, result.Title)
			assert.Equal(t, expected.Body, result.Body)
			assert.Equal(t, expected.Status, result.Status)
			assert.Equal(t, expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, expected.UpdatedAt, result.UpdatedAt)

			if expected.CompletedAt != nil {
				require.NotNil(t, result.CompletedAt)
				assert.Equal(t, *expected.CompletedAt, *result.CompletedAt)
			} else {
				assert.Nil(t, result.CompletedAt)
			}
		})
	}
}

func TestDomainStatusToProtoStatus(t *testing.T) {
	tests := []struct {
		name           string
		domainStatus   todo.TodoStatus
		expectedStatus pb.TodoStatus
	}{
		{
			name:           "converts not started status",
			domainStatus:   todo.TodoStatusNotStarted,
			expectedStatus: pb.TodoStatus_TODO_STATUS_NOT_STARTED,
		},
		{
			name:           "converts in progress status",
			domainStatus:   todo.TodoStatusInProgress,
			expectedStatus: pb.TodoStatus_TODO_STATUS_IN_PROGRESS,
		},
		{
			name:           "converts completed status",
			domainStatus:   todo.TodoStatusCompleted,
			expectedStatus: pb.TodoStatus_TODO_STATUS_COMPLETED,
		},
		{
			name:           "converts unknown status to unspecified",
			domainStatus:   todo.TodoStatus(999), // Invalid status
			expectedStatus: pb.TodoStatus_TODO_STATUS_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domainStatusToProtoStatus(tt.domainStatus)
			assert.Equal(t, tt.expectedStatus, result)
		})
	}
}

func TestParseUUIDFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		validate    func(t *testing.T, result todo.TodoID, err error)
	}{
		{
			name:        "parses valid UUID",
			input:       "550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
			validate: func(t *testing.T, result todo.TodoID, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.String())
			},
		},
		{
			name:        "returns error for invalid UUID format",
			input:       "invalid-uuid",
			expectError: true,
			validate: func(t *testing.T, result todo.TodoID, err error) {
				assert.Error(t, err)
				assert.Equal(t, todo.TodoID{}, result)
			},
		},
		{
			name:        "returns error for empty string",
			input:       "",
			expectError: true,
			validate: func(t *testing.T, result todo.TodoID, err error) {
				assert.Error(t, err)
				assert.Equal(t, todo.TodoID{}, result)
			},
		},
		{
			name:        "returns error for malformed UUID",
			input:       "550e8400-e29b-41d4-a716",
			expectError: true,
			validate: func(t *testing.T, result todo.TodoID, err error) {
				assert.Error(t, err)
				assert.Equal(t, todo.TodoID{}, result)
			},
		},
		{
			name:        "parses UUID with different case",
			input:       "550E8400-E29B-41D4-A716-446655440000",
			expectError: false,
			validate: func(t *testing.T, result todo.TodoID, err error) {
				assert.NoError(t, err)
				// UUID should be normalized to lowercase
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseUUIDFromString(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.validate(t, result, err)
		})
	}
}
