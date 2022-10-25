package storage

import (
	"github.com/sujit-baniya/framework/utils"
	"sync"
	"sync/atomic"
	"time"
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

// Storage interface that is implemented by storage providers
type Storage struct {
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
func New(config ...Config) *Storage {
	// Set default config
	cfg := configDefault(config...)

	// Create storage
	store := &Storage{
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
func (s *Storage) Get(key string) ([]byte, error) {
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
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
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
func (s *Storage) Delete(key string) error {
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
func (s *Storage) Reset() error {
	s.mux.Lock()
	s.db = make(map[string]entry)
	s.mux.Unlock()
	return nil
}

// Close the memory storage
func (s *Storage) Close() error {
	s.done <- struct{}{}
	return nil
}

func (s *Storage) gc() {
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

// Return database client
func (s *Storage) Conn() map[string]entry {
	return s.db
}
