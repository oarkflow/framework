package http

type Json map[string]any

type Response interface {
	String(code int, format string, values ...any) error
	Json(code int, obj any) error
	File(filepath string, compress ...bool) error
	Download(filepath, filename string) error
	Success() ResponseSuccess
	StatusCode() int
	Header(key, value string) Response
	Vary(key string, value ...string)
}

type ResponseSuccess interface {
	String(format string, values ...any) error
	Json(obj any) error
}
