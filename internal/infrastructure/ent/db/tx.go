package db

import (
	"context"
	"errors"

	"github.com/iktakahiro/oniongo/internal/application"
	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/entgen"
)

type key int

const (
	TxKey key = iota
)

// entTransactionManager is the implementation of the TransactionManager interface.
type entTransactionManager struct{}

// NewEntTransactionManager creates a new ent transaction manager.
func NewEntTransactionManager() application.TransactionManager {
	return &entTransactionManager{}
}

// RunInTx runs a function in a transaction.
func (m entTransactionManager) RunInTx(
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
