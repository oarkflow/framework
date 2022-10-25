package route

import (
	"net/http"

	httpcontract "github.com/sujit-baniya/framework/contracts/http"
)

type GroupFunc func(routes Route)

type Engine interface {
	Route
	Run(addr string) error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type Route interface {
	Group(GroupFunc)
	Prefix(addr string) Route
	GlobalMiddleware(...httpcontract.HandlerFunc) Route
	Middleware(...httpcontract.HandlerFunc) Route

	Any(string, ...httpcontract.HandlerFunc)
	Get(string, ...httpcontract.HandlerFunc)
	Post(string, ...httpcontract.HandlerFunc)
	Delete(string, ...httpcontract.HandlerFunc)
	Patch(string, ...httpcontract.HandlerFunc)
	Put(string, ...httpcontract.HandlerFunc)
	Options(string, ...httpcontract.HandlerFunc)

	Static(prefix string, dir string)
	StaticFile(prefix string, file string)
	StaticFS(dir string, fs http.FileSystem)
}
