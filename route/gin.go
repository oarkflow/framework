package route

import (
	"fmt"
	"github.com/sujit-baniya/framework/view"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gookit/color"

	httpcontract "github.com/sujit-baniya/framework/contracts/http"
	"github.com/sujit-baniya/framework/contracts/route"
	"github.com/sujit-baniya/framework/facades"
	"github.com/sujit-baniya/framework/foundation"
	frameworkhttp "github.com/sujit-baniya/framework/http"
)

type Gin struct {
	route.Route
	instance *gin.Engine
	config   frameworkhttp.GinConfig
	view     *view.Engine
}

func NewGin(cfg ...frameworkhttp.GinConfig) route.Engine {
	var config frameworkhttp.GinConfig
	if len(cfg) > 0 {
		config = cfg[0]
	}
	if config.Mode == "" {
		config.Mode = gin.ReleaseMode
	}
	if config.Extension == "" {
		config.Extension = ".html"
	}
	gin.SetMode(config.Mode)
	engine := gin.New()
	if config.Path != "" && config.View == nil {
		config.View = view.New(config.Path, config.Extension)
	}
	return &Gin{instance: engine, view: config.View, Route: NewGinGroup(
		engine.Group("/"),
		"",
		[]httpcontract.HandlerFunc{},
		[]httpcontract.HandlerFunc{},
		config,
		config.View,
	)}
}

func (r *Gin) Run(addr string) error {
	rootApp := foundation.Application{}
	if facades.Config.GetBool("app.debug") && !rootApp.RunningInConsole() {
		routes := r.instance.Routes()
		for _, item := range routes {
			fmt.Printf("%-10s %s\n", item.Method, colonToBracket(item.Path))
		}
	}

	color.Greenln("Listening and serving HTTP on " + addr)

	return r.instance.Run([]string{addr}...)
}

func (r *Gin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.instance.ServeHTTP(w, req)
}

type GinGroup struct {
	instance          gin.IRouter
	originPrefix      string
	originMiddlewares []httpcontract.HandlerFunc
	prefix            string
	middlewares       []httpcontract.HandlerFunc
	globalMiddlewares []httpcontract.HandlerFunc
	config            frameworkhttp.GinConfig
	view              *view.Engine
}

func NewGinGroup(instance gin.IRouter, prefix string, originMiddlewares []httpcontract.HandlerFunc, globalMiddlewares []httpcontract.HandlerFunc, config frameworkhttp.GinConfig, engine *view.Engine) route.Route {
	return &GinGroup{
		instance:          instance,
		originPrefix:      prefix,
		originMiddlewares: originMiddlewares,
		globalMiddlewares: globalMiddlewares,
		config:            config,
		view:              engine,
	}
}

func (r *GinGroup) Group(handler route.GroupFunc) {
	var middlewares []httpcontract.HandlerFunc
	middlewares = append(middlewares, r.originMiddlewares...)
	middlewares = append(middlewares, r.middlewares...)
	r.middlewares = []httpcontract.HandlerFunc{}
	prefix := pathToGinPath(r.originPrefix + "/" + r.prefix)
	r.prefix = ""

	handler(NewGinGroup(r.instance, prefix, middlewares, r.globalMiddlewares, r.config, r.view))
}

func (r *GinGroup) Prefix(addr string) route.Route {
	r.prefix += "/" + addr
	return r
}

func (r *GinGroup) Middleware(handlers ...httpcontract.HandlerFunc) route.Route {
	r.middlewares = append(r.middlewares, handlers...)
	return r
}

func (r *GinGroup) GlobalMiddleware(handlers ...httpcontract.HandlerFunc) route.Route {
	r.globalMiddlewares = append(r.globalMiddlewares, handlers...)
	return r
}

func (r *GinGroup) Any(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []gin.HandlerFunc
	for _, ha := range handlers {
		h = append(h, handlerToGinHandler(ha, r.config, r.view))
	}
	r.getGinRoutesWithMiddlewares().Any(pathToGinPath(relativePath), h...)
}

func (r *GinGroup) Get(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []gin.HandlerFunc
	for _, ha := range handlers {
		h = append(h, handlerToGinHandler(ha, r.config, r.view))
	}
	r.getGinRoutesWithMiddlewares().GET(pathToGinPath(relativePath), h...)
}

func (r *GinGroup) Post(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []gin.HandlerFunc
	for _, ha := range handlers {
		h = append(h, handlerToGinHandler(ha, r.config, r.view))
	}
	r.getGinRoutesWithMiddlewares().POST(pathToGinPath(relativePath), h...)
}

func (r *GinGroup) Delete(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []gin.HandlerFunc
	for _, ha := range handlers {
		h = append(h, handlerToGinHandler(ha, r.config, r.view))
	}
	r.getGinRoutesWithMiddlewares().DELETE(pathToGinPath(relativePath), h...)
}

func (r *GinGroup) Patch(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []gin.HandlerFunc
	for _, ha := range handlers {
		h = append(h, handlerToGinHandler(ha, r.config, r.view))
	}
	r.getGinRoutesWithMiddlewares().PATCH(pathToGinPath(relativePath), h...)
}

func (r *GinGroup) Put(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []gin.HandlerFunc
	for _, ha := range handlers {
		h = append(h, handlerToGinHandler(ha, r.config, r.view))
	}
	r.getGinRoutesWithMiddlewares().PUT(pathToGinPath(relativePath), h...)
}

func (r *GinGroup) Options(relativePath string, handlers ...httpcontract.HandlerFunc) {
	var h []gin.HandlerFunc
	for _, ha := range handlers {
		h = append(h, handlerToGinHandler(ha, r.config, r.view))
	}
	r.getGinRoutesWithMiddlewares().OPTIONS(pathToGinPath(relativePath), h...)
}

func (r *GinGroup) Static(relativePath, root string) {
	r.getGinRoutesWithMiddlewares().Static(pathToGinPath(relativePath), root)
}

func (r *GinGroup) StaticFile(relativePath, filepath string) {
	r.getGinRoutesWithMiddlewares().StaticFile(pathToGinPath(relativePath), filepath)
}

func (r *GinGroup) StaticFS(relativePath string, fs http.FileSystem) {
	r.getGinRoutesWithMiddlewares().StaticFS(pathToGinPath(relativePath), fs)
}

func (r *GinGroup) getGinRoutesWithMiddlewares() gin.IRoutes {
	prefix := pathToGinPath(r.originPrefix + "/" + r.prefix)
	r.prefix = ""
	ginGroup := r.instance.Group(prefix)

	var middlewares []gin.HandlerFunc
	ginOriginMiddlewares := middlewaresToGinHandlers(r.originMiddlewares, r.config, r.view)
	ginMiddlewares := middlewaresToGinHandlers(r.middlewares, r.config, r.view)
	ginGlobalMiddlewares := middlewaresToGinHandlers(r.globalMiddlewares, r.config, r.view)
	middlewares = append(middlewares, ginOriginMiddlewares...)
	middlewares = append(middlewares, ginMiddlewares...)
	middlewares = append(middlewares, ginGlobalMiddlewares...)
	// middlewares = addDebugLog(middlewares)
	r.middlewares = []httpcontract.HandlerFunc{}

	if len(middlewares) > 0 {
		return ginGroup.Use(middlewares...)
	} else {
		return ginGroup
	}
}

func pathToGinPath(relativePath string) string {
	return bracketToColon(path.Clean(relativePath))
}

func middlewaresToGinHandlers(middlewares []httpcontract.HandlerFunc, config frameworkhttp.GinConfig, engine *view.Engine) []gin.HandlerFunc {
	var ginHandlers []gin.HandlerFunc
	for _, item := range middlewares {
		ginHandlers = append(ginHandlers, middlewareToGinHandler(item, config, engine))
	}

	return ginHandlers
}

func handlerToGinHandler(handler httpcontract.HandlerFunc, config frameworkhttp.GinConfig, engine *view.Engine) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		handler(frameworkhttp.NewGinContext(ginCtx, config, engine))
	}
}

func middlewareToGinHandler(handler httpcontract.HandlerFunc, config frameworkhttp.GinConfig, engine *view.Engine) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		handler(frameworkhttp.NewGinContext(ginCtx, config, engine))
	}
}

func addDebugLog(middlewares []gin.HandlerFunc) []gin.HandlerFunc {
	logFormatter := func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			param.Latency = param.Latency - param.Latency%time.Second
		}
		return fmt.Sprintf("[HTTP] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}

	if facades.Config.GetBool("app.debug") {
		middlewares = append(middlewares, gin.LoggerWithFormatter(logFormatter))
	}

	return middlewares
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

func bracketToColon(relativePath string) string {
	compileRegex := regexp.MustCompile("\\{(.*?)\\}")
	matchArr := compileRegex.FindAllStringSubmatch(relativePath, -1)

	for _, item := range matchArr {
		relativePath = strings.ReplaceAll(relativePath, item[0], ":"+item[1])
	}

	return relativePath
}
