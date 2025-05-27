package todohandler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	v1 "github.com/iktakahiro/oniongo/internal/api/grpc/gen/oniongo/v1"
	"github.com/iktakahiro/oniongo/internal/application/todoapp"
	"github.com/samber/do"
)

// CreateTodoHandler handles CreateTodo requests
type createTodoHandler struct {
	useCase todoapp.CreateTodoUseCase
}

func newCreateTodoHandler(i *do.Injector) (*createTodoHandler, error) {
	createTodoUseCase, err := do.Invoke[todoapp.CreateTodoUseCase](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke create todo use case: %w", err)
	}
	return &createTodoHandler{useCase: createTodoUseCase}, nil
}

func (h createTodoHandler) CreateTodo(
	ctx context.Context,
	req *connect.Request[v1.CreateTodoRequest],
) (*connect.Response[v1.CreateTodoResponse], error) {
	// Extract body value if present
	body := ""
	if req.Msg.Body != nil {
		body = *req.Msg.Body
	}

	// Create use case request
	useCaseReq := todoapp.CreateTodoRequest{
		Title: req.Msg.Title,
		Body:  body,
	}

	// Execute use case
	if err := h.useCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.CreateTodoResponse{}), nil
}
