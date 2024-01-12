package facades

import (
	"github.com/oarkflow/framework/storage"
)

var Storage *storage.Storage

func init() {
	Storage = storage.NewStorage("")
}
