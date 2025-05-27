package todohandler

import (
	"github.com/google/uuid"
	"github.com/iktakahiro/oniongo/internal/domain/todo"
	pb "github.com/iktakahiro/oniongo/internal/infrastructure/grpc/gen/oniongo/v1"
)

// domainTodoToProto converts a domain Todo to a protobuf Todo
func domainTodoToProto(domainTodo *todo.Todo) *pb.Todo {
	pbTodo := &pb.Todo{
		Id:        domainTodo.ID().String(),
		Title:     domainTodo.Title(),
		Body:      domainTodo.Body(),
		Status:    domainStatusToProtoStatus(domainTodo.Status()),
		CreatedAt: domainTodo.CreatedAt().Unix(),
		UpdatedAt: domainTodo.UpdatedAt().Unix(),
	}

	if completedAt := domainTodo.CompletedAt(); completedAt != nil {
		timestamp := completedAt.Unix()
		pbTodo.CompletedAt = &timestamp
	}

	return pbTodo
}

// domainStatusToProtoStatus converts a domain TodoStatus to a protobuf TodoStatus
func domainStatusToProtoStatus(domainStatus todo.TodoStatus) pb.TodoStatus {
	switch domainStatus {
	case todo.TodoStatusNotStarted:
		return pb.TodoStatus_TODO_STATUS_NOT_STARTED
	case todo.TodoStatusInProgress:
		return pb.TodoStatus_TODO_STATUS_IN_PROGRESS
	case todo.TodoStatusCompleted:
		return pb.TodoStatus_TODO_STATUS_COMPLETED
	default:
		return pb.TodoStatus_TODO_STATUS_UNSPECIFIED
	}
}

// parseUUIDFromString parses a UUID string and returns a TodoID
func parseUUIDFromString(idStr string) (todo.TodoID, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return todo.TodoID{}, err
	}
	return todo.TodoID(id), nil
}
