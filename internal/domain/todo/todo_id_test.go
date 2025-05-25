package todo

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewTodoID(t *testing.T) {
	// When
	id1 := NewTodoID()
	id2 := NewTodoID()

	// Then
	require.NotEqual(t, TodoID{}, id1)
	require.NotEqual(t, TodoID{}, id2)
	require.NotEqual(t, id1, id2) // Each call should generate a unique ID
	require.NotEmpty(t, id1.String())
	require.NotEmpty(t, id2.String())
}

func TestTodoID_String(t *testing.T) {
	tests := []struct {
		name string
		id   TodoID
	}{
		{
			name: "valid todo id",
			id:   NewTodoID(),
		},
		{
			name: "zero todo id",
			id:   TodoID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			result := tt.id.String()

			// Then
			require.NotNil(t, result)
			if tt.id != (TodoID{}) {
				require.NotEmpty(t, result)
				// Verify it's a valid UUID string format
				_, err := uuid.Parse(result)
				require.NoError(t, err)
			}
		})
	}
}

func TestTodoID_UUID(t *testing.T) {
	tests := []struct {
		name string
		id   TodoID
	}{
		{
			name: "valid todo id",
			id:   NewTodoID(),
		},
		{
			name: "zero todo id",
			id:   TodoID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			result := tt.id.UUID()

			// Then
			require.NotNil(t, result)
			// Verify the conversion is consistent
			require.Equal(t, uuid.UUID(tt.id), result)
		})
	}
}

func TestNewTodoIDFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid uuid string",
			input:       "550e8400-e29b-41d4-a716-446655440000",
			expectError: false,
		},
		{
			name:        "valid uuid v7 string",
			input:       NewTodoID().String(),
			expectError: false,
		},
		{
			name:        "invalid uuid string",
			input:       "invalid-uuid",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "malformed uuid",
			input:       "550e8400-e29b-41d4-a716",
			expectError: true,
		},
		{
			name:        "uuid with invalid characters",
			input:       "550e8400-e29b-41d4-a716-44665544000g",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			result, err := NewTodoIDFromString(tt.input)

			// Then
			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, TodoID{}, result)
			} else {
				require.NoError(t, err)
				require.NotEqual(t, TodoID{}, result)
				require.Equal(t, tt.input, result.String())
			}
		})
	}
}
