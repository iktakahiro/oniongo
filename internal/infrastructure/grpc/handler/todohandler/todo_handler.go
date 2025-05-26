package todohandler

import (
	"context"

	"connectrpc.com/connect"
	"github.com/iktakahiro/oniongo/internal/application/todoapp"
	v1 "github.com/iktakahiro/oniongo/internal/infrastructure/grpc/gen/oniongo/v1"
	v1connect "github.com/iktakahiro/oniongo/internal/infrastructure/grpc/gen/oniongo/v1/oniongov1connect"
	"github.com/samber/do"
)

type todoHandler struct {
	createTodoUseCase   todoapp.CreateTodoUseCase
	getTodoUseCase      todoapp.GetTodoUseCase
	getTodosUseCase     todoapp.GetTodosUseCase
	updateTodoUseCase   todoapp.UpdateTodoUseCase
	startTodoUseCase    todoapp.StartTodoUseCase
	completeTodoUseCase todoapp.CompleteTodoUseCase
	deleteTodoUseCase   todoapp.DeleteTodoUseCase
}

func NewTodoServiceHandler(i *do.Injector) (v1connect.TodoServiceHandler, error) {
	createTodoUseCase := do.MustInvoke[todoapp.CreateTodoUseCase](i)
	getTodoUseCase := do.MustInvoke[todoapp.GetTodoUseCase](i)
	getTodosUseCase := do.MustInvoke[todoapp.GetTodosUseCase](i)
	updateTodoUseCase := do.MustInvoke[todoapp.UpdateTodoUseCase](i)
	startTodoUseCase := do.MustInvoke[todoapp.StartTodoUseCase](i)
	completeTodoUseCase := do.MustInvoke[todoapp.CompleteTodoUseCase](i)
	deleteTodoUseCase := do.MustInvoke[todoapp.DeleteTodoUseCase](i)
	return &todoHandler{
		createTodoUseCase:   createTodoUseCase,
		getTodoUseCase:      getTodoUseCase,
		getTodosUseCase:     getTodosUseCase,
		updateTodoUseCase:   updateTodoUseCase,
		startTodoUseCase:    startTodoUseCase,
		completeTodoUseCase: completeTodoUseCase,
		deleteTodoUseCase:   deleteTodoUseCase,
	}, nil
}

func (h *todoHandler) CreateTodo(
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
	if err := h.createTodoUseCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.CreateTodoResponse{}), nil
}

func (h *todoHandler) GetTodo(
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
	domainTodo, err := h.getTodoUseCase.Execute(ctx, useCaseReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Convert to protobuf and return response
	pbTodo := domainTodoToProto(domainTodo)
	return connect.NewResponse(&v1.GetTodoResponse{
		Todo: pbTodo,
	}), nil
}

func (h *todoHandler) GetTodos(
	ctx context.Context,
	req *connect.Request[v1.GetTodosRequest],
) (*connect.Response[v1.GetTodosResponse], error) {
	// Create use case request
	useCaseReq := todoapp.GetTodosRequest{}

	// Execute use case
	domainTodos, err := h.getTodosUseCase.Execute(ctx, useCaseReq)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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

func (h *todoHandler) UpdateTodo(
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
	if err := h.updateTodoUseCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.UpdateTodoResponse{}), nil
}

func (h *todoHandler) StartTodo(
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
	if err := h.startTodoUseCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.StartTodoResponse{}), nil
}

func (h *todoHandler) CompleteTodo(
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
	if err := h.completeTodoUseCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.CompleteTodoResponse{}), nil
}

func (h *todoHandler) DeleteTodo(
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
	if err := h.deleteTodoUseCase.Execute(ctx, useCaseReq); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Return response
	return connect.NewResponse(&v1.DeleteTodoResponse{}), nil
}
