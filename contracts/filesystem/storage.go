package filesystem

import (
	"context"
	"time"
)

type Storage interface {
	Driver
	Disk(disk string) Driver
}

type Driver interface {
	WithContext(ctx context.Context) Driver
	Put(file string, content []byte) error
	PutFile(path string, source File) (string, error)
	PutFileAs(path string, source File, name string) (string, error)
	Get(file string) ([]byte, error)
	Size(file string) (int64, error)
	Path(file string) string
	Exists(file string) bool
	Missing(file string) bool
	// Url Download(path string)
	Url(file string) string
	TemporaryUrl(file string, time time.Time) (string, error)
	Copy(oldFile, newFile string) error
	Move(oldFile, newFile string) error
	Delete(file ...string) error
	Files(path string) ([]string, error)
	AllFiles(path string) ([]string, error)
	Directories(path string) ([]string, error)
	AllDirectories(path string) ([]string, error)
	MakeDirectory(directory string) error
	DeleteDirectory(directory string) error
}

type File interface {
	Disk(disk string) File
	File() string
	Store(path string) (string, error)
	StoreAs(path string, name string) (string, error)
	GetClientOriginalName() string
	GetClientOriginalExtension() string
	HashName(path ...string) string
	Extension() (string, error)
}

type Option func(options *Options)

type Options struct {
	Name string
}
