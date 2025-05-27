package todohandler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	v1 "github.com/iktakahiro/oniongo/internal/api/grpc/gen/oniongo/v1"
	"github.com/iktakahiro/oniongo/internal/application/todoapp"
	"github.com/samber/do"
)

// StartTodoHandler handles StartTodo requests
type startTodoHandler struct {
	useCase todoapp.StartTodoUseCase
}

func newStartTodoHandler(i *do.Injector) (*startTodoHandler, error) {
	startTodoUseCase, err := do.Invoke[todoapp.StartTodoUseCase](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke start todo use case: %w", err)
	}
	return &startTodoHandler{useCase: startTodoUseCase}, nil
}

func (h startTodoHandler) StartTodo(
	ctx context.Context,
	req *connect.Request[v1.StartTodoRequest],
) (*connect.Response[v1.StartTodoResponse], error) {
	// Parse todo ID
	todoID, err := parseUUIDFromString(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Create use case request
	useCaseReq := todoapp.StartTodoRequest{
		ID: todoID,
	}

	// Execute use case
	if err := h.useCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.StartTodoResponse{}), nil
}
