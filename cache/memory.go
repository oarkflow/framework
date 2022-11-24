package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sujit-baniya/framework/contracts/cache"
	contractStorage "github.com/sujit-baniya/framework/contracts/storage"
	"github.com/sujit-baniya/framework/storage"
)

type Memory struct {
	Prefix string
	Client contractStorage.Storage
}

func NewMemoryCache(prefix string) cache.Store {
	return &Memory{
		Prefix: prefix,
		Client: storage.New(),
	}
}

// WithContext Retrieve an item from the cache by key.
func (r *Memory) WithContext(ctx context.Context) cache.Store {
	return r
}

// Get Retrieve an item from the cache by key.
func (r *Memory) Get(key string, def interface{}) interface{} {
	val, err := r.Client.Get(r.Prefix + key)
	if err != nil {
		switch s := def.(type) {
		case func() interface{}:
			return s()
		default:
			return def
		}
	}

	return val
}

func (r *Memory) GetBool(key string, def bool) bool {
	res := r.Get(key, def)
	switch val := res.(type) {
	case string:
		switch val {
		case "0", "false":
			return false
		case "1", "true":
			return true
		}
	case []byte:
		if len(val) == 0 {
			return false
		}
		v := string(val)
		switch v {
		case "0", "false":
			return false
		case "1", "true":
			return true
		}
	}
	return res.(bool)
}

func (r *Memory) GetInt(key string, def int) int {
	res := r.Get(key, def)
	if val, ok := res.(string); ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return def
		}

		return i
	}

	return res.(int)
}

func (r *Memory) GetString(key string, def string) string {
	data := r.Get(key, def)
	switch d := data.(type) {
	case []byte:
		return string(d)
	default:
		return fmt.Sprintf("%v", d)
	}
}

// Has Check an item exists in the cache.
func (r *Memory) Has(key string) bool {
	value, err := r.Client.Get(r.Prefix + key)
	if err != nil || value == nil {
		return false
	}

	return true
}

// Put Store an item in the cache for a given number of seconds.
func (r *Memory) Put(key string, value interface{}, seconds time.Duration) error {
	err := r.Client.Set(r.Prefix+key, []byte(fmt.Sprint(value)), seconds)
	if err != nil {
		return err
	}

	return nil
}

// Pull Retrieve an item from the cache and delete it.
func (r *Memory) Pull(key string, def interface{}) interface{} {
	val, err := r.Client.Get(r.Prefix + key)
	if err != nil {
		return def
	}
	err = r.Client.Delete(r.Prefix + key)

	if err != nil {
		return def
	}

	return val
}

// Add Store an item in the cache if the key does not exist.
func (r *Memory) Add(key string, value interface{}, seconds time.Duration) bool {
	err := r.Client.Set(r.Prefix+key, []byte(fmt.Sprint(value)), seconds)
	if err != nil {
		return false
	}

	return true
}

// Remember Get an item from the cache, or execute the given Closure and store the result.
func (r *Memory) Remember(key string, ttl time.Duration, callback func() interface{}) (interface{}, error) {
	val := r.Get(key, nil)

	if val != nil {
		return val, nil
	}

	val = callback()

	if err := r.Put(key, val, ttl); err != nil {
		return nil, err
	}

	return val, nil
}

// RememberForever Get an item from the cache, or execute the given Closure and store the result forever.
func (r *Memory) RememberForever(key string, callback func() interface{}) (interface{}, error) {
	val := r.Get(key, nil)

	if val != nil {
		return val, nil
	}

	val = callback()

	if err := r.Put(key, val, 0); err != nil {
		return nil, err
	}

	return val, nil
}

// Forever Store an item in the cache indefinitely.
func (r *Memory) Forever(key string, value interface{}) bool {
	if err := r.Put(key, value, 0); err != nil {
		return false
	}

	return true
}

// Forget Remove an item from the cache.
func (r *Memory) Forget(key string) bool {
	err := r.Client.Delete(r.Prefix + key)

	if err != nil {
		return false
	}

	return true
}

// Flush Remove all items from the cache.
func (r *Memory) Flush() bool {
	err := r.Client.Reset()

	if err != nil {
		return false
	}

	return true
}
