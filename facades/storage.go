package facades

import (
	"github.com/sujit-baniya/framework/contracts/storage"
	store "github.com/sujit-baniya/framework/storage"
)

var Memory storage.Storage

func init() {
	Memory = store.New()
}
