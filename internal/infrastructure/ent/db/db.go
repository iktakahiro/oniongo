package db

import (
	"sync"

	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/entgen"
)

var (
	clientInstance *entgen.Client
	clientOnce     sync.Once
	clientErr      error
)

// GetClient returns a singleton instance of the database client
func GetClient() (*entgen.Client, error) {
	clientOnce.Do(func() {
		clientInstance, clientErr = entgen.Open(
			"sqlite3",
			"file:ent?mode=memory&cache=shared&_fk=1",
		)
	})
	return clientInstance, clientErr
}
