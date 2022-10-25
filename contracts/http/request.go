package http

import (
	"mime/multipart"
	"net/http"
)

type Request interface {
	Origin() *http.Request
}

type File interface {
	SaveFile(name string, dst string) error
	File(name string) (*multipart.FileHeader, error)
}
