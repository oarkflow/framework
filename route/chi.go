package route

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	httpcontract "github.com/sujit-baniya/framework/contracts/http"
	"github.com/sujit-baniya/framework/contracts/route"
	frameworkhttp "github.com/sujit-baniya/framework/http"
)

type Chi struct {
	route.Route
	instance chi.Router
}

func NewChi() route.Engine {
	engine := chi.NewRouter()

	return &Chi{instance: engine, Route: NewChiGroup(
		engine,
		"",
		[]httpcontract.HandlerFunc{},
		[]httpcontract.HandlerFunc{},
	)}
}

func (r *Chi) Run(addr string) error {
	// @TODO - Implement
	/*rootApp := foundation.Application{}
	if facades.Config.GetBool("app.debug") && !rootApp.RunningInConsole() {
		routes := r.instance.Routes()
		for _, item := range routes {
			fmt.Printf("%-10s %s\n", item, colonToBracket(item.Path))
		}
	}*/

	// color.Greenln("Listening and serving HTTP on " + addr)

	return http.ListenAndServe(addr, r.instance)
}

func (r *Chi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.instance.ServeHTTP(w, req)
}

type ChiGroup struct {
	instance          chi.Router
	originPrefix      string
	originMiddlewares []httpcontract.HandlerFunc
	prefix            string
	middlewares       []httpcontract.HandlerFunc
	globalMiddlewares []httpcontract.HandlerFunc
}

func NewChiGroup(instance chi.Router, prefix string, originMiddlewares []httpcontract.HandlerFunc, globalMiddlewares []httpcontract.HandlerFunc) route.Route {
	return &ChiGroup{
		instance:          instance,
		originPrefix:      prefix,
		originMiddlewares: originMiddlewares,
		globalMiddlewares: globalMiddlewares,
	}
}

func (r *ChiGroup) Group(handler route.GroupFunc) {
	var middlewares []httpcontract.HandlerFunc
	middlewares = append(middlewares, r.originMiddlewares...)
	middlewares = append(middlewares, r.middlewares...)
	r.middlewares = []httpcontract.HandlerFunc{}
	prefix := pathToChiPath(r.originPrefix + "/" + r.prefix)
	r.prefix = ""

	handler(NewChiGroup(r.instance, prefix, middlewares, r.globalMiddlewares))
}

func (r *ChiGroup) Prefix(addr string) route.Route {
	r.prefix += "/" + addr
	return r
}

func (r *ChiGroup) Middleware(handlers ...httpcontract.HandlerFunc) route.Route {
	r.middlewares = append(r.middlewares, handlers...)
	return r
}

func (r *ChiGroup) GlobalMiddleware(handlers ...httpcontract.HandlerFunc) route.Route {
	r.globalMiddlewares = append(r.globalMiddlewares, handlers...)
	return r
}

func (r *ChiGroup) Any(relativePath string, handlers ...httpcontract.HandlerFunc) {
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1])
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha))
	}
	m := middlewaresToChiHandlers(middlewares)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Get(pathToChiPath(relativePath), handler)
	ro.With(m...).Post(pathToChiPath(relativePath), handler)
	ro.With(m...).Put(pathToChiPath(relativePath), handler)
	ro.With(m...).Delete(pathToChiPath(relativePath), handler)
	ro.With(m...).Patch(pathToChiPath(relativePath), handler)
	ro.With(m...).Options(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Get(relativePath string, handlers ...httpcontract.HandlerFunc) {
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1])
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha))
	}
	m := middlewaresToChiHandlers(middlewares)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Get(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Post(relativePath string, handlers ...httpcontract.HandlerFunc) {
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1])
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha))
	}
	m := middlewaresToChiHandlers(middlewares)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Post(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Delete(relativePath string, handlers ...httpcontract.HandlerFunc) {
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1])
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha))
	}
	m := middlewaresToChiHandlers(middlewares)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Delete(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Patch(relativePath string, handlers ...httpcontract.HandlerFunc) {
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1])
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha))
	}
	m := middlewaresToChiHandlers(middlewares)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Patch(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Put(relativePath string, handlers ...httpcontract.HandlerFunc) {
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1])
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha))
	}
	m := middlewaresToChiHandlers(middlewares)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Put(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Options(relativePath string, handlers ...httpcontract.HandlerFunc) {
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1])
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha))
	}
	m := middlewaresToChiHandlers(middlewares)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Options(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Static(relativePath, root string) {
	//@TODO - Implement
}

func (r *ChiGroup) StaticFile(relativePath, filepath string) {
	//@TODO - Implement
}

func (r *ChiGroup) StaticFS(relativePath string, fs http.FileSystem) {
	//@TODO - Implement
}

func (r *ChiGroup) getChiRoutesWithMiddlewares() chi.Router {
	var middlewares []func(handler http.Handler) http.Handler
	ginOriginMiddlewares := middlewaresToChiHandlers(r.originMiddlewares)
	ginMiddlewares := middlewaresToChiHandlers(r.middlewares)
	ginGlobalMiddlewares := middlewaresToChiHandlers(r.globalMiddlewares)
	middlewares = append(middlewares, ginOriginMiddlewares...)
	middlewares = append(middlewares, ginMiddlewares...)
	middlewares = append(middlewares, ginGlobalMiddlewares...)
	r.middlewares = []httpcontract.HandlerFunc{}
	return r.instance
}

func pathToChiPath(relativePath string) string {
	return mergeSlashForPath(relativePath)
}

func middlewaresToChiHandlers(middlewares []httpcontract.HandlerFunc) []func(handler http.Handler) http.Handler {
	var ginHandlers []func(handler http.Handler) http.Handler
	for _, item := range middlewares {
		ginHandlers = append(ginHandlers, middlewareToChiHandler(item))
	}

	return ginHandlers
}

func handlerToChiHandler(handler httpcontract.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, req *http.Request) {
		handler(frameworkhttp.NewChiContext(req, response))
	}
}

func middlewareToChiHandler(handler httpcontract.HandlerFunc) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			handler(frameworkhttp.NewChiContext(request, writer))
		})
	}
}
