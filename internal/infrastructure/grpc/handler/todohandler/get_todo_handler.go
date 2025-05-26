package todohandler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/iktakahiro/oniongo/internal/application/todoapp"
	v1 "github.com/iktakahiro/oniongo/internal/infrastructure/grpc/gen/oniongo/v1"
	"github.com/samber/do"
)

// GetTodoHandler handles GetTodo requests
type getTodoHandler struct {
	useCase todoapp.GetTodoUseCase
}

func newGetTodoHandler(i *do.Injector) (*getTodoHandler, error) {
	getTodoUseCase, err := do.Invoke[todoapp.GetTodoUseCase](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke get todo use case: %w", err)
	}
	return &getTodoHandler{useCase: getTodoUseCase}, nil
}

func (h getTodoHandler) GetTodo(
	ctx context.Context,
	req *connect.Request[v1.GetTodoRequest],
) (*connect.Response[v1.GetTodoResponse], error) {
	// Parse todo ID
	todoID, err := parseUUIDFromString(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Create use case request
	useCaseReq := todoapp.GetTodoRequest{
		ID: todoID,
	}

	// Execute use case
	domainTodo, err := h.useCase.Execute(ctx, useCaseReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Convert to protobuf and return response
	pbTodo := domainTodoToProto(domainTodo)
	return connect.NewResponse(&v1.GetTodoResponse{
		Todo: pbTodo,
	}), nil
}
