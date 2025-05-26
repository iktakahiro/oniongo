package todohandler

import (
	v1connect "github.com/iktakahiro/oniongo/internal/infrastructure/grpc/gen/oniongo/v1/oniongov1connect"
	"github.com/samber/do"
)

// todoServiceHandler combines all individual handlers to implement TodoServiceHandler
type todoServiceHandler struct {
	*createTodoHandler
	*getTodoHandler
	*getTodosHandler
	*updateTodoHandler
	*startTodoHandler
	*completeTodoHandler
	*deleteTodoHandler
}

// NewTodoServiceHandler creates a new TodoServiceHandler using composition
func NewTodoServiceHandler(i *do.Injector) (v1connect.TodoServiceHandler, error) {
	createHandler, err := newCreateTodoHandler(i)
	if err != nil {
		return nil, err
	}
	getHandler, err := newGetTodoHandler(i)
	if err != nil {
		return nil, err
	}
	getTodosHandler, err := newGetTodosHandler(i)
	if err != nil {
		return nil, err
	}
	updateHandler, err := newUpdateTodoHandler(i)
	if err != nil {
		return nil, err
	}
	startHandler, err := newStartTodoHandler(i)
	if err != nil {
		return nil, err
	}
	completeHandler, err := newCompleteTodoHandler(i)
	if err != nil {
		return nil, err
	}
	deleteHandler, err := newDeleteTodoHandler(i)
	if err != nil {
		return nil, err
	}

	return &todoServiceHandler{
		createTodoHandler:   createHandler,
		getTodoHandler:      getHandler,
		getTodosHandler:     getTodosHandler,
		updateTodoHandler:   updateHandler,
		startTodoHandler:    startHandler,
		completeTodoHandler: completeHandler,
		deleteTodoHandler:   deleteHandler,
	}, nil
}
