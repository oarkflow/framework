package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gookit/color"

	"github.com/oarkflow/framework/cache/stores/memory"
	"github.com/oarkflow/framework/cache/stores/redis"
	"github.com/oarkflow/framework/contracts/cache"
	"github.com/oarkflow/framework/facades"
)

type Cache struct {
	cache.Store
	Name       string
	ctx        context.Context
	connection string
	instances  map[string]cache.Store
	mu         *sync.RWMutex
}

var (
	defaultCache *Cache
)

func Connection(name string) cache.Store {
	defaultCache.mu.Lock()
	defer defaultCache.mu.Unlock()
	defaultConnection := facades.Config.GetString("cache.default")
	if name == "" {
		name = defaultConnection
	}

	if defaultCache.instances == nil {
		defaultCache.instances = make(map[string]cache.Store)
	}
	if st, exist := defaultCache.instances[name]; exist {
		return st
	}
	g, err := NewDriver(defaultCache.ctx, name)
	if err != nil {
		color.Redln(fmt.Sprintf("[Filesystem] Init connection error, %v", err))

		return nil
	}
	if g == nil {
		return nil
	}

	defaultCache.instances[name] = g

	if name == defaultConnection {
		defaultCache.Store = g
	}

	return g
}

func NewCache(name string) *Cache {
	str := &Cache{Name: name, mu: &sync.RWMutex{}, instances: make(map[string]cache.Store)}
	str.Connection(name)
	return str
}

func (s *Cache) Connection(name string) cache.Store {
	s.mu.Lock()
	defer s.mu.Unlock()
	defaultConnection := facades.Config.GetString("cache.default")
	if name == "" {
		name = defaultConnection
	}

	s.connection = name
	if s.instances == nil {
		s.instances = make(map[string]cache.Store)
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
		s.Store = g
	}

	return g
}

func NewDriver(ctx context.Context, name string) (cache.Store, error) {
	driver := facades.Config.GetString(fmt.Sprintf("cache.stores.%s.driver", name))
	switch driver {
	case "redis":
		return NewRedisCache(ctx, name)
	default:
		return memory.New(""), nil
	}
}

func NewRedisCache(ctx context.Context, name string) (cache.Store, error) {
	config := facades.Config.Get(fmt.Sprintf("cache.stores.%s", name))
	switch config := config.(type) {
	case map[string]any:
		var cfg redis.Config
		bt, _ := json.Marshal(config)
		json.Unmarshal(bt, &cfg)
		cfg.Context = ctx
		return redis.New(cfg)
	default:
		cfg := redis.Config{
			Host:    "localhost",
			Port:    "6379",
			DB:      0,
			Context: ctx,
		}
		return redis.New(cfg)
	}
}
