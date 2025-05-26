// Package uow provides a interface as Unit of Work pattern.
package uow

import "context"

// TransactionRunner is the interface for running operations within a transaction.
type TransactionRunner interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
