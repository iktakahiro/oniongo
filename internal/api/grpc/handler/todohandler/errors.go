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

	// Check error types using errors.As
	var notFoundErr *domainTodo.NotFoundError
	if errors.As(err, &notFoundErr) {
		return connect.NewError(connect.CodeNotFound, err)
	}

	var validationErr *domainTodo.ValidationError
	if errors.As(err, &validationErr) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	var stateErr *domainTodo.StateError
	if errors.As(err, &stateErr) {
		return connect.NewError(connect.CodeFailedPrecondition, err)
	}

	// Default to internal error
	return connect.NewError(connect.CodeInternal, err)
}
