// Package uow provides a interface as Unit of Work pattern.
package uow

import "context"

// TransactionManager is the interface for the transaction manager.
type TransactionManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
