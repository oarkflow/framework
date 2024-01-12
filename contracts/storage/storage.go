package storage

import (
	"time"
)

// Storage interface for communicating with different database/key-value
// providers
type Storage interface {
	Reset() error
	Set(key string, value []byte, exp time.Duration) (err error)
	Delete(key string) (err error)
	Get(key string) (value []byte, err error)
	Connection(name string) Storage
	Iterate(fn func(key []byte, value []byte))
	IterateByPrefix(prefix []byte, limit uint64, fn func(key []byte, value []byte)) uint64
	IterateByPrefixFrom(prefix []byte, from []byte, limit uint64, fn func(key []byte, value []byte)) uint64
	DeleteByPrefix(prefix []byte)
	KeysByPrefixCount(prefix []byte) uint64
	Close() error
}
