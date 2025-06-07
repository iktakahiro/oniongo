package todohandler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	v1 "github.com/iktakahiro/oniongo/internal/api/grpc/gen/oniongo/v1"
	"github.com/iktakahiro/oniongo/internal/application/todoapp"
	"github.com/samber/do"
)

// GetTodosHandler handles GetTodos requests
type getTodosHandler struct {
	useCase todoapp.GetTodosUseCase
}

func newGetTodosHandler(i *do.Injector) (*getTodosHandler, error) {
	getTodosUseCase, err := do.Invoke[todoapp.GetTodosUseCase](i)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke get todos use case: %w", err)
	}
	return &getTodosHandler{useCase: getTodosUseCase}, nil
}

func (h getTodosHandler) GetTodos(
	ctx context.Context,
	req *connect.Request[v1.GetTodosRequest],
) (*connect.Response[v1.GetTodosResponse], error) {
	// Create use case request
	useCaseReq := todoapp.GetTodosRequest{}

	// Execute use case
	domainTodos, err := h.useCase.Execute(ctx, useCaseReq)
	if err != nil {
		return nil, toConnectError(err)
	}

	// Convert to protobuf
	pbTodos := make([]*v1.Todo, len(domainTodos))
	for i, domainTodo := range domainTodos {
		pbTodos[i] = domainTodoToProto(domainTodo)
	}

	// Return response
	return connect.NewResponse(&v1.GetTodosResponse{
		Todos: pbTodos,
	}), nil
}
