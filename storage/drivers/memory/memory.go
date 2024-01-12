package memory

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/oarkflow/framework/contracts/storage"
	"github.com/oarkflow/framework/utils"
)

// Config defines the config for storage.
type Config struct {
	// Time before deleting expired keys
	//
	// Default is 10 * time.Second
	GCInterval time.Duration
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	GCInterval: 10 * time.Second,
}

// configDefault is a helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if int(cfg.GCInterval.Seconds()) <= 0 {
		cfg.GCInterval = ConfigDefault.GCInterval
	}
	return cfg
}

// Memory interface that is implemented by storage providers
type Memory struct {
	mux        sync.RWMutex
	db         map[string]entry
	gcInterval time.Duration
	done       chan struct{}
}

type entry struct {
	data []byte
	// max value is 4294967295 -> Sun Feb 07 2106 06:28:15 GMT+0000
	expiry uint32
}

// New creates a new memory storage
func New(config ...Config) *Memory {
	// Set default config
	cfg := configDefault(config...)

	// Create storage
	store := &Memory{
		db:         make(map[string]entry),
		gcInterval: cfg.GCInterval,
		done:       make(chan struct{}),
	}

	// Start garbage collector
	utils.StartTimeStampUpdater()
	go store.gc()

	return store
}

// Get value by key
func (s *Memory) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	s.mux.RLock()
	v, ok := s.db[key]
	s.mux.RUnlock()
	if !ok || v.expiry != 0 && v.expiry <= atomic.LoadUint32(&utils.Timestamp) {
		return nil, nil
	}

	return v.data, nil
}

// Set key with value
func (s *Memory) Set(key string, val []byte, exp time.Duration) error {
	// Ain't Nobody Got Time For That
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}

	var expire uint32
	if exp != 0 {
		expire = uint32(exp.Seconds()) + atomic.LoadUint32(&utils.Timestamp)
	}

	s.mux.Lock()
	s.db[key] = entry{val, expire}
	s.mux.Unlock()
	return nil
}

// Delete key by key
func (s *Memory) Delete(key string) error {
	// Ain't Nobody Got Time For That
	if len(key) <= 0 {
		return nil
	}
	s.mux.Lock()
	delete(s.db, key)
	s.mux.Unlock()
	return nil
}

// Reset all keys
func (s *Memory) Reset() error {
	s.mux.Lock()
	s.db = make(map[string]entry)
	s.mux.Unlock()
	return nil
}

// Close the memory storage
func (s *Memory) Close() error {
	s.done <- struct{}{}
	return nil
}

func (s *Memory) gc() {
	ticker := time.NewTicker(s.gcInterval)
	defer ticker.Stop()
	var expired []string

	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			expired = expired[:0]
			s.mux.RLock()
			for id, v := range s.db {
				if v.expiry != 0 && v.expiry < atomic.LoadUint32(&utils.Timestamp) {
					expired = append(expired, id)
				}
			}
			s.mux.RUnlock()
			s.mux.Lock()
			for i := range expired {
				delete(s.db, expired[i])
			}
			s.mux.Unlock()
		}
	}
}

// Conn database client
func (s *Memory) Conn() map[string]entry {
	return s.db
}

func (s *Memory) Connection(name string) storage.Storage {
	//TODO implement me
	panic("implement me")
}

func (s *Memory) Iterate(fn func(key []byte, value []byte)) {
	//TODO implement me
	panic("implement me")
}

func (s *Memory) IterateByPrefix(prefix []byte, limit uint64, fn func(key []byte, value []byte)) uint64 {
	//TODO implement me
	panic("implement me")
}

func (s *Memory) IterateByPrefixFrom(prefix []byte, from []byte, limit uint64, fn func(key []byte, value []byte)) uint64 {
	//TODO implement me
	panic("implement me")
}

func (s *Memory) DeleteByPrefix(prefix []byte) {
	//TODO implement me
	panic("implement me")
}

func (s *Memory) KeysByPrefixCount(prefix []byte) uint64 {
	//TODO implement me
	panic("implement me")
}
