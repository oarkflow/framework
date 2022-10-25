package route

import (
	"fmt"
	"github.com/sujit-baniya/chi"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/foundation"
	"net/http"
	"path"
	"strings"

	httpcontract "github.com/sujit-baniya/framework/contracts/http"
	"github.com/sujit-baniya/framework/contracts/route"
	frameworkhttp "github.com/sujit-baniya/framework/http"
)

type Chi struct {
	route.Route
	instance chi.Router
	config   frameworkhttp.ChiConfig
}

func NewChi(config ...frameworkhttp.ChiConfig) route.Engine {
	var cfg frameworkhttp.ChiConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	engine := chi.NewRouter(chi.Config{
		NotFoundHandler:         cfg.NotFoundHandler,
		MethodNotAllowedHandler: cfg.MethodNotAllowedHandler,
	})

	return &Chi{instance: engine, Route: NewChiGroup(
		engine,
		"/",
		[]httpcontract.HandlerFunc{},
		[]httpcontract.HandlerFunc{},
		cfg,
	)}
}

func (r *Chi) Run(addr string) error {
	// @TODO - Implement
	rootApp := foundation.Application{}
	if facades.Config.GetBool("app.debug") && !rootApp.RunningInConsole() {
		routes := r.instance.Routes()
		for _, item := range routes {
			for method, _ := range item.Handlers {
				fmt.Printf("%-10s %s\n", method, colonToBracket(item.Pattern))
			}
		}
	}

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
	config            frameworkhttp.ChiConfig
}

func NewChiGroup(instance chi.Router, prefix string, originMiddlewares []httpcontract.HandlerFunc, globalMiddlewares []httpcontract.HandlerFunc, config frameworkhttp.ChiConfig) route.Route {
	return &ChiGroup{
		instance:          instance,
		originPrefix:      prefix,
		originMiddlewares: originMiddlewares,
		globalMiddlewares: globalMiddlewares,
		config:            config,
	}
}

func (r *ChiGroup) Group(handler route.GroupFunc) {
	var middlewares []httpcontract.HandlerFunc
	middlewares = append(middlewares, r.originMiddlewares...)
	middlewares = append(middlewares, r.middlewares...)
	r.middlewares = []httpcontract.HandlerFunc{}
	prefix := pathToChiPath(r.originPrefix + "/" + r.prefix)
	r.prefix = ""

	handler(NewChiGroup(r.instance, prefix, middlewares, r.globalMiddlewares, r.config))
}

func (r *ChiGroup) Prefix(addr string) route.Route {
	if strings.HasPrefix(addr, "/") {
		r.prefix = addr
	} else {
		r.prefix += "/" + addr
	}
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
	relativePath = r.originPrefix + "/" + r.prefix + "/" + relativePath
	r.prefix = ""
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1], r.config)
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha, r.config))
	}
	m := middlewaresToChiHandlers(middlewares, r.config)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Get(pathToChiPath(relativePath), handler)
	ro.With(m...).Post(pathToChiPath(relativePath), handler)
	ro.With(m...).Put(pathToChiPath(relativePath), handler)
	ro.With(m...).Delete(pathToChiPath(relativePath), handler)
	ro.With(m...).Patch(pathToChiPath(relativePath), handler)
	ro.With(m...).Options(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Get(relativePath string, handlers ...httpcontract.HandlerFunc) {

	relativePath = r.originPrefix + "/" + r.prefix + "/" + relativePath
	r.prefix = ""
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1], r.config)
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha, r.config))
	}
	m := middlewaresToChiHandlers(middlewares, r.config)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Get(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Post(relativePath string, handlers ...httpcontract.HandlerFunc) {

	relativePath = r.originPrefix + "/" + r.prefix + "/" + relativePath
	r.prefix = ""
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1], r.config)
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha, r.config))
	}
	m := middlewaresToChiHandlers(middlewares, r.config)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Post(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Delete(relativePath string, handlers ...httpcontract.HandlerFunc) {

	relativePath = r.originPrefix + "/" + r.prefix + "/" + relativePath
	r.prefix = ""
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1], r.config)
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha, r.config))
	}
	m := middlewaresToChiHandlers(middlewares, r.config)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Delete(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Patch(relativePath string, handlers ...httpcontract.HandlerFunc) {

	relativePath = r.originPrefix + "/" + r.prefix + "/" + relativePath
	r.prefix = ""
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1], r.config)
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha, r.config))
	}
	m := middlewaresToChiHandlers(middlewares, r.config)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Patch(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Put(relativePath string, handlers ...httpcontract.HandlerFunc) {

	relativePath = r.originPrefix + "/" + r.prefix + "/" + relativePath
	r.prefix = ""
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1], r.config)
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha, r.config))
	}
	m := middlewaresToChiHandlers(middlewares, r.config)
	ro := r.getChiRoutesWithMiddlewares()
	ro.With(m...).Put(pathToChiPath(relativePath), handler)
}

func (r *ChiGroup) Options(relativePath string, handlers ...httpcontract.HandlerFunc) {

	relativePath = r.originPrefix + "/" + r.prefix + "/" + relativePath
	r.prefix = ""
	middlewares := handlers[0 : len(handlers)-1]
	handler := handlerToChiHandler(handlers[len(handlers)-1], r.config)
	var h []http.HandlerFunc
	for _, ha := range middlewares {
		h = append(h, handlerToChiHandler(ha, r.config))
	}
	m := middlewaresToChiHandlers(middlewares, r.config)
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
	ginOriginMiddlewares := middlewaresToChiHandlers(r.originMiddlewares, r.config)
	ginMiddlewares := middlewaresToChiHandlers(r.middlewares, r.config)
	ginGlobalMiddlewares := middlewaresToChiHandlers(r.globalMiddlewares, r.config)
	middlewares = append(middlewares, ginOriginMiddlewares...)
	middlewares = append(middlewares, ginMiddlewares...)
	middlewares = append(middlewares, ginGlobalMiddlewares...)
	r.middlewares = []httpcontract.HandlerFunc{}
	if len(middlewares) > 0 {
		r.instance.With(middlewares...)
	}
	return r.instance
}

func pathToChiPath(relativePath string) string {
	return path.Clean(relativePath)
}

func middlewaresToChiHandlers(middlewares []httpcontract.HandlerFunc, config frameworkhttp.ChiConfig) []func(handler http.Handler) http.Handler {
	var ginHandlers []func(handler http.Handler) http.Handler
	for _, item := range middlewares {
		ginHandlers = append(ginHandlers, middlewareToChiHandler(item, config))
	}

	return ginHandlers
}

func handlerToChiHandler(handler httpcontract.HandlerFunc, config frameworkhttp.ChiConfig) http.HandlerFunc {
	return func(response http.ResponseWriter, req *http.Request) {
		handler(frameworkhttp.NewChiContext(req, response, config))
	}
}

func middlewareToChiHandler(handler httpcontract.HandlerFunc, config frameworkhttp.ChiConfig) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			handler(frameworkhttp.NewChiContext(request, writer, config, next))
		})
	}
}

func colonToBracket(relativePath string) string {
	arr := strings.Split(relativePath, "/")
	var newArr []string
	for _, item := range arr {
		if strings.HasPrefix(item, ":") {
			item = "{" + strings.ReplaceAll(item, ":", "") + "}"
		}
		newArr = append(newArr, item)
	}

	return strings.Join(newArr, "/")
}
