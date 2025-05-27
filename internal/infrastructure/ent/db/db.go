package db

import (
	"context"
	"sync"

	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/entgen"
	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/entgen/migrate"
	_ "github.com/mattn/go-sqlite3"
)

var (
	clientInstance *entgen.Client
	clientOnce     sync.Once
	clientErr      error
)

// GetClient returns a singleton instance of the database client
func GetClient() (*entgen.Client, error) {
	if clientInstance != nil {
		return clientInstance, nil
	}

	clientOnce.Do(func() {
		clientInstance, clientErr = entgen.Open(
			"sqlite3",
			"file:db/dev.db?_fk=1",
		)
	})
	return clientInstance, clientErr
}

func Migrate() error {
	db, err := GetClient()
	if err != nil {
		return err
	}

	return db.Schema.Create(context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
}
