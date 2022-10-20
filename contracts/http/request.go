package http

import (
	"mime/multipart"
	"net/http"
)

type Request interface {
	Origin() *http.Request
	Response() Response

	//Validate(ctx *gin.Context, request FormRequest) []error
}

type File interface {
	Store(dst string) error
	File() *multipart.FileHeader
}
