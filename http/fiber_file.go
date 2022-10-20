package http

import (
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
)

type FiberFile struct {
	instance *fiber.Ctx
	file     *multipart.FileHeader
}

func (f *FiberFile) Store(dst string) error {
	return f.instance.SaveFile(f.file, dst)
}

func (f *FiberFile) File() *multipart.FileHeader {
	return f.file
}
