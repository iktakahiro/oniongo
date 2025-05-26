package todohandler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/iktakahiro/oniongo/internal/application/todoapp"
	v1 "github.com/iktakahiro/oniongo/internal/infrastructure/grpc/gen/oniongo/v1"
	"github.com/samber/do"
)

// UpdateTodoHandler handles UpdateTodo requests
type updateTodoHandler struct {
	useCase todoapp.UpdateTodoUseCase
}

func newUpdateTodoHandler(i *do.Injector) (*updateTodoHandler, error) {
	updateTodoUseCase, err := do.Invoke[todoapp.UpdateTodoUseCase](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke update todo use case: %w", err)
	}
	return &updateTodoHandler{useCase: updateTodoUseCase}, nil
}

func (h updateTodoHandler) UpdateTodo(
	ctx context.Context,
	req *connect.Request[v1.UpdateTodoRequest],
) (*connect.Response[v1.UpdateTodoResponse], error) {
	// Parse todo ID
	todoID, err := parseUUIDFromString(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Extract body value if present
	body := ""
	if req.Msg.Body != nil {
		body = *req.Msg.Body
	}

	// Create use case request
	useCaseReq := todoapp.UpdateTodoRequest{
		ID:    todoID,
		Title: req.Msg.Title,
		Body:  body,
	}

	// Execute use case
	if err := h.useCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.UpdateTodoResponse{}), nil
}
