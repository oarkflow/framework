package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/gookit/color"

	"github.com/oarkflow/framework/contracts/storage"
	"github.com/oarkflow/framework/facades"
	"github.com/oarkflow/framework/storage/drivers/badger"
	"github.com/oarkflow/framework/storage/drivers/memory"
)

type Storage struct {
	storage.Storage
	Name       string
	ctx        context.Context
	connection string
	instances  map[string]storage.Storage
	mu         *sync.RWMutex
}

var (
	defaultStorage *Storage
)

func Connection(name string) storage.Storage {
	defaultStorage.mu.Lock()
	defer defaultStorage.mu.Unlock()
	defaultConnection := facades.Config.GetString("storage.default")
	if name == "" {
		name = defaultConnection
	}

	defaultStorage.connection = name
	if defaultStorage.instances == nil {
		defaultStorage.instances = make(map[string]storage.Storage)
	}

	if st, exist := defaultStorage.instances[name]; exist {
		return st
	}

	g, err := NewDriver(defaultStorage.ctx, name)
	if err != nil {
		color.Redln(fmt.Sprintf("[Filesystem] Init connection error, %v", err))

		return nil
	}
	if g == nil {
		return nil
	}

	defaultStorage.instances[name] = g

	if name == defaultConnection {
		defaultStorage.Storage = g
	}

	return defaultStorage
}

func NewStorage(name string) *Storage {
	str := &Storage{Name: name, mu: &sync.RWMutex{}}
	str.Connection(name)
	return str
}

func (s *Storage) Connection(name string) storage.Storage {
	s.mu.Lock()
	defer s.mu.Unlock()
	defaultConnection := facades.Config.GetString("storage.default")
	if name == "" {
		name = defaultConnection
	}

	s.connection = name
	if s.instances == nil {
		s.instances = make(map[string]storage.Storage)
	}

	if _, exist := s.instances[name]; exist {
		return s
	}

	g, err := NewDriver(s.ctx, name)
	if err != nil {
		color.Redln(fmt.Sprintf("[Filesystem] Init connection error, %v", err))

		return nil
	}
	if g == nil {
		return nil
	}

	s.instances[name] = g

	if name == defaultConnection {
		s.Storage = g
	}

	return s
}

func NewDriver(ctx context.Context, name string) (storage.Storage, error) {
	driver := facades.Config.GetString(fmt.Sprintf("storage.drivers.%s.driver", name))
	switch driver {
	case "badger":
		return NewBadgerDB(name)
	default:
		return memory.New(), nil
	}
}

func NewBadgerDB(name string) (storage.Storage, error) {
	config := facades.Config.Get(fmt.Sprintf("storage.drivers.%s", name))
	dbPath := "/tmp/badger"
	switch config := config.(type) {
	case map[string]any:
		if path, e := config["database"]; e {
			dbPath = path.(string)
		}
		opts := badger.DefaultOptions(dbPath)
		if mem, e := config["in_memory"]; e {
			opts.InMemory = mem.(bool)
		}
		if mem, e := config["index_cache"]; e {
			opts.IndexCacheSize = mem.(int64)
		}
		if mem, e := config["block_cache"]; e {
			opts.BlockCacheSize = mem.(int64)
		}
		return badger.New(dbPath, opts)
	default:
		opts := badger.DefaultOptions(dbPath)
		return badger.New(dbPath, opts)
	}
}
