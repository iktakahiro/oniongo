package todohandler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	v1 "github.com/iktakahiro/oniongo/internal/api/grpc/gen/oniongo/v1"
	"github.com/iktakahiro/oniongo/internal/application/todoapp"
	"github.com/samber/do"
)

// CompleteTodoHandler handles CompleteTodo requests
type completeTodoHandler struct {
	useCase todoapp.CompleteTodoUseCase
}

func newCompleteTodoHandler(i *do.Injector) (*completeTodoHandler, error) {
	completeTodoUseCase, err := do.Invoke[todoapp.CompleteTodoUseCase](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke complete todo use case: %w", err)
	}
	return &completeTodoHandler{useCase: completeTodoUseCase}, nil
}

func (h completeTodoHandler) CompleteTodo(
	ctx context.Context,
	req *connect.Request[v1.CompleteTodoRequest],
) (*connect.Response[v1.CompleteTodoResponse], error) {
	// Parse todo ID
	todoID, err := parseUUIDFromString(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Create use case request
	useCaseReq := todoapp.CompleteTodoRequest{
		ID: todoID,
	}

	// Execute use case
	if err := h.useCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.CompleteTodoResponse{}), nil
}
