package todohandler

import (
	"errors"

	"connectrpc.com/connect"
	domainTodo "github.com/iktakahiro/oniongo/internal/domain/todo"
)

// toConnectError converts domain errors to appropriate Connect error codes
func toConnectError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domainTodo.ErrNotFound):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.Is(err, domainTodo.ErrTitleRequired):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, domainTodo.ErrInvalidStateTransition):
		return connect.NewError(connect.CodeFailedPrecondition, err)
	case errors.Is(err, domainTodo.ErrAlreadyCompleted):
		return connect.NewError(connect.CodeFailedPrecondition, err)
	default:
		return connect.NewError(connect.CodeInternal, err)
	}
}