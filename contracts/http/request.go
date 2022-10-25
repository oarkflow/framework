package http

import (
	"mime/multipart"
)

type Request interface {
	Origin() any
}

type File interface {
	SaveFile(name string, dst string) error
	File(name string) (*multipart.FileHeader, error)
}
