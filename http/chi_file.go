package http

import (
	"mime/multipart"
	"net/http"
)

type ChiFile struct {
	request *http.Request
	file    *multipart.FileHeader
}

func (f *ChiFile) Store(dst string) error {
	// @TODO - implement
	return nil // f.Req.MultipartForm.File
}

func (f *ChiFile) File() *multipart.FileHeader {
	return f.file
}
