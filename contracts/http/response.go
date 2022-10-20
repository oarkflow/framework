package http

type Json map[string]interface{}

type Response interface {
	String(code int, format string, values ...interface{}) error
	Json(code int, obj interface{}) error
	File(filepath string, compress ...bool) error
	Download(filepath, filename string) error
	Success() ResponseSuccess
	Header(key, value string) Response
}

type ResponseSuccess interface {
	String(format string, values ...interface{}) error
	Json(obj interface{}) error
}
