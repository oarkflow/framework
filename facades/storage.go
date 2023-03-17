package facades

import (
	"github.com/oarkflow/framework/contracts/storage"
	store "github.com/oarkflow/framework/storage"
)

var Memory storage.Storage

func init() {
	Memory = store.New()
}
