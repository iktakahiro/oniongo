package db

import (
	"context"
	"errors"

	"github.com/iktakahiro/oniongo/internal/application/uow"
	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/entgen"
	"github.com/samber/do"
)

type key int

const (
	TxKey key = iota
)

// entTransactionRunner is the implementation of the TransactionRunner interface.
type entTransactionRunner struct{}

// NewEntTransactionRunner creates a new ent transaction runner.
func NewEntTransactionRunner(i *do.Injector) (uow.TransactionRunner, error) {
	return &entTransactionRunner{}, nil
}

// RunInTx runs a function in a transaction.
func (r entTransactionRunner) RunInTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	tx, err := client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var done bool

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	// set the transaction to the context
	// this is used to get the transaction in the repository
	// e.g.
	//   tx, _ := GetTx(ctx)
	//   tx.Create(ctx, &entgen.Todo{
	// 	  Title: "test",
	// 	  Body:  "test",
	//   })
	ctx = context.WithValue(ctx, TxKey, tx)

	if err := fn(ctx); err != nil {
		return err
	}
	done = true

	return tx.Commit()
}

// GetTx returns the transaction from the context.
func GetTx(ctx context.Context) (*entgen.Tx, error) {
	tx, ok := ctx.Value(TxKey).(*entgen.Tx)
	if !ok {
		return nil, errors.New("tx not found")
	}
	return tx, nil
}
