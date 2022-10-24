package route

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"path"

	httpcontract "github.com/sujit-baniya/framework/contracts/http"
	"github.com/sujit-baniya/framework/contracts/route"
	frameworkhttp "github.com/sujit-baniya/framework/http"
)

type Fiber struct {
	route.Route
	instance *fiber.App
}

func NewFiber(config ...fiber.Config) route.Engine {
	engine := fiber.New(config...)
	return &Fiber{instance: engine, Route: NewFiberGroup(
		engine,
		"/",
		[]httpcontract.HandlerFunc{},
		[]httpcontract.HandlerFunc{},
	)}
}

func (r *Fiber) Run(addr string) error {
	return r.instance.Listen(addr)
}

func (r *Fiber) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

type FiberGroup struct {
	instance          fiber.Router
	originPrefix      string
	originMiddlewares []httpcontract.HandlerFunc
	prefix            string
	middlewares       []httpcontract.HandlerFunc
	globalMiddlewares []httpcontract.HandlerFunc
}

func NewFiberGroup(instance fiber.Router, prefix string, originMiddlewares []httpcontract.HandlerFunc, globalMiddlewares []httpcontract.HandlerFunc) route.Route {
	return &FiberGroup{
		instance:          instance,
		originPrefix:      prefix,
		originMiddlewares: originMiddlewares,
		globalMiddlewares: globalMiddlewares,
	}
}

func (r *FiberGroup) Group(handler route.GroupFunc) {
	var middlewares []httpcontract.HandlerFunc
	middlewares = append(middlewares, r.originMiddlewares...)
	middlewares = append(middlewares, r.middlewares...)
	r.middlewares = []httpcontract.HandlerFunc{}
	prefix := pathToFiberPath(r.originPrefix + "/" + r.prefix)
	r.prefix = ""

	handler(NewFiberGroup(r.instance, prefix, middlewares, r.globalMiddlewares))
}

func (r *FiberGroup) Prefix(addr string) route.Route {
	r.prefix += "/" + addr
	return r
}

func (r *FiberGroup) Middleware(handlers ...httpcontract.HandlerFunc) route.Route {
	r.middlewares = append(r.middlewares, handlers...)
	return r
}

func (r *FiberGroup) GlobalMiddleware(handlers ...httpcontract.HandlerFunc) route.Route {
	r.globalMiddlewares = append(r.globalMiddlewares, handlers...)
	return r
}

func (r *FiberGroup) Any(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []fiber.Handler
	for _, ha := range handlers {
		h = append(h, handlerToFiberHandler(ha))
	}
	r.getFiberRoutesWithMiddlewares().All(pathToFiberPath(relativePath), h...)
}

func (r *FiberGroup) Get(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []fiber.Handler
	for _, ha := range handlers {
		h = append(h, handlerToFiberHandler(ha))
	}
	r.getFiberRoutesWithMiddlewares().Get(pathToFiberPath(relativePath), h...)
}

func (r *FiberGroup) Post(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []fiber.Handler
	for _, ha := range handlers {
		h = append(h, handlerToFiberHandler(ha))
	}
	r.getFiberRoutesWithMiddlewares().Post(pathToFiberPath(relativePath), h...)
}

func (r *FiberGroup) Delete(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []fiber.Handler
	for _, ha := range handlers {
		h = append(h, handlerToFiberHandler(ha))
	}
	r.getFiberRoutesWithMiddlewares().Delete(pathToFiberPath(relativePath), h...)
}

func (r *FiberGroup) Patch(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []fiber.Handler
	for _, ha := range handlers {
		h = append(h, handlerToFiberHandler(ha))
	}
	r.getFiberRoutesWithMiddlewares().Patch(pathToFiberPath(relativePath), h...)
}

func (r *FiberGroup) Put(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []fiber.Handler
	for _, ha := range handlers {
		h = append(h, handlerToFiberHandler(ha))
	}
	r.getFiberRoutesWithMiddlewares().Put(pathToFiberPath(relativePath), h...)
}

func (r *FiberGroup) Options(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []fiber.Handler
	for _, ha := range handlers {
		h = append(h, handlerToFiberHandler(ha))
	}
	r.getFiberRoutesWithMiddlewares().Options(pathToFiberPath(relativePath), h...)
}

func (r *FiberGroup) Static(relativePath, root string) {
	r.getFiberRoutesWithMiddlewares().Static(pathToFiberPath(relativePath), root)
}

func (r *FiberGroup) StaticFile(relativePath, filepath string) {

}

func (r *FiberGroup) StaticFS(relativePath string, fs http.FileSystem) {

}

func (r *FiberGroup) getFiberRoutesWithMiddlewares() fiber.Router {
	prefix := pathToFiberPath(r.originPrefix + "/" + r.prefix)
	r.prefix = ""

	var middlewares []fiber.Handler
	ginOriginMiddlewares := middlewaresToFiberHandlers(r.originMiddlewares)
	ginMiddlewares := middlewaresToFiberHandlers(r.middlewares)
	ginGlobalMiddlewares := middlewaresToFiberHandlers(r.globalMiddlewares)
	middlewares = append(middlewares, ginOriginMiddlewares...)
	middlewares = append(middlewares, ginMiddlewares...)
	middlewares = append(middlewares, ginGlobalMiddlewares...)
	r.middlewares = []httpcontract.HandlerFunc{}

	if len(middlewares) > 0 {
		return r.instance.Group(prefix, middlewares...)
	} else {
		return r.instance
	}
}

func pathToFiberPath(relativePath string) string {
	return bracketToColon(path.Clean(relativePath))
}

func middlewaresToFiberHandlers(middlewares []httpcontract.HandlerFunc) []fiber.Handler {
	var ginHandlers []fiber.Handler
	for _, item := range middlewares {
		ginHandlers = append(ginHandlers, middlewareToFiberHandler(item))
	}

	return ginHandlers
}

func handlerToFiberHandler(handler httpcontract.HandlerFunc) fiber.Handler {
	return func(ginCtx *fiber.Ctx) error {
		return handler(frameworkhttp.NewFiberContext(ginCtx))
	}
}

func middlewareToFiberHandler(handler httpcontract.HandlerFunc) fiber.Handler {
	return func(ginCtx *fiber.Ctx) error {
		return handler(frameworkhttp.NewFiberContext(ginCtx))
	}
}
