package todohandler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	v1 "github.com/iktakahiro/oniongo/internal/api/grpc/gen/oniongo/v1"
	"github.com/iktakahiro/oniongo/internal/application/todoapp"
	"github.com/samber/do"
)

// DeleteTodoHandler handles DeleteTodo requests
type deleteTodoHandler struct {
	useCase todoapp.DeleteTodoUseCase
}

func newDeleteTodoHandler(i *do.Injector) (*deleteTodoHandler, error) {
	deleteTodoUseCase, err := do.Invoke[todoapp.DeleteTodoUseCase](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke delete todo use case: %w", err)
	}
	return &deleteTodoHandler{useCase: deleteTodoUseCase}, nil
}

func (h deleteTodoHandler) DeleteTodo(
	ctx context.Context,
	req *connect.Request[v1.DeleteTodoRequest],
) (*connect.Response[v1.DeleteTodoResponse], error) {
	// Parse todo ID
	todoID, err := parseUUIDFromString(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Create use case request
	useCaseReq := todoapp.DeleteTodoRequest{
		ID: todoID,
	}

	// Execute use case
	if err := h.useCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.DeleteTodoResponse{}), nil
}
