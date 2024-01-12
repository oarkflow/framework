package badger

import (
	"time"

	"github.com/dgraph-io/badger/v4"

	"github.com/oarkflow/framework/contracts/storage"
	"github.com/oarkflow/framework/facades"
)

// Badger implements wrapper for badger database
type Badger struct {
	db *badger.DB
}

func DefaultOptions(path string) badger.Options {
	return badger.DefaultOptions(path)
}

// New returns new instance of badger wrapper
func New(storageDir string, opts badger.Options) (*Badger, error) {
	storage := &Badger{}
	opts.ZSTDCompressionLevel = 10
	opts.SyncWrites = true
	opts.Dir = storageDir
	opts.ValueDir = storageDir
	var err error
	storage.db, err = badger.Open(opts)
	if err != nil {
		return nil, err
	}

	go storage.runStorageGC()

	return storage, nil
}

func (storage *Badger) Connection(name string) storage.Storage {
	return facades.Storage.Connection(name)
}

// Close properly closes badger database
func (storage *Badger) Close() error {
	return storage.db.Close()
}

// Set adds a key-value pair to the database
func (storage *Badger) Set(key string, value []byte, exp time.Duration) (err error) {
	return storage.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), value)
		return err
	})
}

// Delete deletes a key
func (storage *Badger) Delete(key string) (err error) {
	return storage.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
}

// Reset deletes a key
func (storage *Badger) Reset() (err error) {
	return storage.db.DropAll()
}

// Get returns value by key
func (storage *Badger) Get(key string) (value []byte, err error) {
	err = storage.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(value)
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// Iterate iterates over all keys
func (storage *Badger) Iterate(fn func(key []byte, value []byte)) {
	storage.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.KeyCopy(nil)
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			fn(k, v)
		}
		return nil
	})
}

// IterateByPrefix iterates over keys with prefix
func (storage *Badger) IterateByPrefix(prefix []byte, limit uint64, fn func(key []byte, value []byte)) uint64 {
	var totalIterated uint64
	storage.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix) && ((limit > 0 && totalIterated < limit) || limit <= 0); it.Next() {
			item := it.Item()
			k := item.KeyCopy(nil)
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			fn(k, v)
			totalIterated++
		}
		return nil
	})

	return totalIterated
}

func (storage *Badger) KeysByPrefixCount(prefix []byte) uint64 {
	var count uint64
	storage.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			count++
		}

		return nil
	})

	return count
}

// DeleteByPrefix iterates over keys with prefix
func (storage *Badger) DeleteByPrefix(prefix []byte) {
	deleteKeys := func(keysForDelete [][]byte) error {
		if err := storage.db.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000
	keysForDeleteBunches := make([][][]byte, 0)
	keysForDelete := make([][]byte, 0, collectSize)
	keysCollected := 0

	// создать банчи и удалять банчами после итератора же ну
	storage.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				keysForDeleteBunches = append(keysForDeleteBunches, keysForDelete)
				keysForDelete = make([][]byte, 0, collectSize)
				keysCollected = 0
			}
		}
		if keysCollected > 0 {
			keysForDeleteBunches = append(keysForDeleteBunches, keysForDelete)
		}

		return nil
	})

	for _, keys := range keysForDeleteBunches {
		deleteKeys(keys)
	}
}

// IterateByPrefixFrom iterates over keys with prefix
func (storage *Badger) IterateByPrefixFrom(prefix []byte, from []byte, limit uint64, fn func(key []byte, value []byte)) uint64 {
	var totalIterated uint64
	storage.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(from); it.ValidForPrefix(prefix) && ((limit > 0 && totalIterated < limit) || limit <= 0); it.Next() {
			item := it.Item()
			k := item.KeyCopy(nil)
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			fn(k, v)
			totalIterated++
		}
		return nil
	})

	return totalIterated
}

func (storage *Badger) runStorageGC() {
	timer := time.NewTicker(10 * time.Minute)
	for range timer.C {
		storage.storageGC()
	}
}

func (storage *Badger) storageGC() {
again:
	err := storage.db.RunValueLogGC(0.5)
	if err == nil {
		goto again
	}
}
