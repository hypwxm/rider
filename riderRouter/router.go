//package router
//路由注册模块，注册http请求（Get,Post,Delete,Option...via）
package riderRouter

import (
	"net/http"
	"sync"
	"rider/riderServer"
	"path/filepath"
	"log"
	"strings"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

type BaseHandler interface {
	ServeHTTP (w http.ResponseWriter, r *http.Request)
}


//根路由或者子路由需要实现的方法
type BaseRouter interface {
	ANY(path string, router *Router) *Router
	POST(path string, router *Router)
	GetServer() *riderServer.HttpServer
	SetServer() *riderServer.HttpServer
	registerHandler(method string, path string, router *Router)
	//ServerFile(path string, fileRoot string)
	GET(path string, router *Router)
	HEAD(path string, router *Router)
	OPTIONS(path string, router *Router)
	PUT(path string, router *Router)
	PATCH(path string, router *Router)
	DELETE(path string, router *Router)
	//HiJack(path string, router *Router)
	//WebSocket(path string, router *Router)
	Any(path string, router *Router)
	RegisterRouter(routeMethod string, path string, router *Router)
}

//根路由配置
type Routers struct {
	routers           map[string]*Router
	mux               *sync.Mutex
	server            *riderServer.HttpServer
	Middleware        []HandlerFunc
	PathMiddleware	  map[string][]HandlerFunc
	registeredRouters RegisteredRouters
}

//子路由配置
type Router struct {
	SubRouter  RegisteredRouters
	mux        *sync.Mutex
	Handler    HandlerFunc
	Middleware []HandlerFunc
	PathMiddleware	  map[string][]HandlerFunc
	Method     string
	FullPath   string
	server     *riderServer.HttpServer
	rootPath   string
	currentPath string
	rootMethod string
}

type handlerRouter map[string]*Router

type RegisteredRouters map[string]handlerRouter

var RegisteredRoutersMux = new(sync.Mutex)

var registeredRouters = make(RegisteredRouters)

func (rr RegisteredRouters) NewHandle() handlerRouter {
	return make(map[string]*Router)
}

func (rr RegisteredRouters) NewPath(path string) {
	if rr[path] == nil {
		rr[path] = make(handlerRouter)
	}
}

func (rr RegisteredRouters) MatchMethodAndPath(method string, path string) int {
	return MatchMethodAndPath(method, path, rr)
}

//给RegisteredRouters添加新的的处理路由
func (rr RegisteredRouters) NewRouter(method string, pattern string, router *Router) {
	//先要验证路由的唯一性
	//判断全局注册的路由和请求方法的存在和存在的状态0，1，2
	registerStatus := MatchMethodAndPath(method, pattern, rr)
	if registerStatus == 0 || registerStatus == 1 {
		RegisteredRoutersMux.Lock()
		rr.NewPath(pattern)
		rr[pattern][method] = router
		RegisteredRoutersMux.Unlock()
	} else if registerStatus == 2 {
		log.Fatalln(method + " " + pattern + " 已注册，请勿重复注册同一个http方法和请求路径")
		return
	}
}


//根路由初始化方法
func NewRouters() *Routers {
	return &Routers{
		mux:               new(sync.Mutex),
		registeredRouters: registeredRouters,
		routers: make(map[string]*Router),
		PathMiddleware: make(map[string][]HandlerFunc),
	}
}

//子路由初始化方法
func NewRouter(handler HandlerFunc) *Router {
	return &Router{
		SubRouter: make(RegisteredRouters),
		Handler:   handler,
		mux:       new(sync.Mutex),
		PathMiddleware: make(map[string][]HandlerFunc),
	}
}

func (r *Routers) Run() {
	for path, handler := range registeredRouters {
		r.server.ServerMux.Handle(path, handler)
	}
}

func (r *Routers) GetRegisteredRouters() RegisteredRouters {
	return registeredRouters
}



func (h handlerRouter) ServeHTTP (w http.ResponseWriter, req *http.Request) {
	//请求的路径
	reqPath := req.URL.Path
	//请求的http方法
	reqMethod := req.Method
	if strings.ToUpper(req.Method) == http.MethodOptions {
		allows := allow(reqPath)
		w.Header().Set("allow", allows)
		return
	}

	var (
		finalRouter *Router
		anyRoute *Router
	)

	//如果方法里面有匹配的处理，则会按照改方法处理，
	//如果不存在匹配，但是存在ANY，则会在未找到其他匹配方式的情况下，用ANY来处理该请求
	//若最后没找到任何可以处理的方式，则返回404
	for method, router := range h {
		if method == "ANY" {
			anyRoute = router
		}

		if method != "ANY" {
			if method == reqMethod {
				router.doHandlers(w, req)
				return
			}
			finalRouter = router
		}
	}

	if anyRoute != nil {
		anyRoute.doHandlers(w, req)
		return
	}

	//当没有任何可匹配的请求响应时，用全局的server（所有的router保存的都是全局的server）来响应404
	finalRouter.server.Error(w, "not found", http.StatusNotFound)
}

//请求方式匹配后将执行该路由注册的handlers和middleware
func (r *Router) doHandlers(w http.ResponseWriter, req *http.Request) {
	//未注册任何中间件和处理方式，直接返回404
	if len(r.PathMiddleware[r.currentPath]) == 0 && r.Handler == nil {
		r.server.Error(w, "not found", http.StatusNotFound)
		return
	}

	//先处理中间件函数（如果存在）
	if len(r.PathMiddleware[r.currentPath]) > 0 {
		for _, handle := range r.PathMiddleware[r.currentPath] {
			handle(w, req)
		}
	}

	//处理完中间件，处理最终函数（如果存在）
	if r.Handler != nil {
		r.Handler(w, req)
	}
}

//获取根路由内部的服务来源
func (r *Routers) GetServer() *riderServer.HttpServer {
	return r.server
}

//将服务注入根路由
func (r *Routers) SetServer(server *riderServer.HttpServer) {
	r.server = server
}

//实现registerHandler
func (r *Routers) registerHandler(method string, path string, router *Router) {

}

//添加中间件(将app服务的中间件加进来)
func (r *Routers) AppendMiddleware(path string, appMiddleware ...HandlerFunc) {
	if r.PathMiddleware == nil {
		r.PathMiddleware = make(map[string][]HandlerFunc)
	}
	r.mux.Lock()
	r.PathMiddleware[path] = append(appMiddleware, r.PathMiddleware[path][:]...)
	r.mux.Unlock()
}

//调用使用添加中间件(将app服务的中间件加进来)
func (r *Router) AddMiddleware(middleware ...HandlerFunc) {
	r.Middleware = append(r.Middleware, middleware...)
}

//当请求为options时，返回该路径所支持的请求方式
func allow(path string) string {
	var allows = ""
	for _path, handlerMethodMap := range registeredRouters {
		if _path == path {
			for _method, _ := range handlerMethodMap {
				if _method != "OPTIONS" {
					if allows == "" {
						allows += _method
					} else {
						allows += "," + _method
					}
				}
			}
		}
	}
	return allows
}

//获取所有已注册的路由
func GetRegisteredRouters () map[string]handlerRouter {
	return registeredRouters
}

//查询routers根路由内部是否已经存在要注册的路由和方法，如果路径和方法都存在了就返回错误（重复注册）,
//如果路径存在，请求方式还未注册，在这路径上在注册一个对应的请求方式
//返回参数： 0：路径和方法都还未注册过
//			1：路径已有注册其他方法，但是方法还没注册
//			2：路径和方法都已经注册，不能再次注册
func MatchMethodAndPath(method string, path string, registeredRouters RegisteredRouters) int {
	var matchInt int
	for _path, registeredMethod := range registeredRouters {
		if _path == path {
			for _method, _ := range registeredMethod {
				if _method == method {
					//log.Fatalln(_method + " " + path + " 已注册，请勿重复注册同一个http方法和请求路径")
					matchInt = 2
					break
				}
				matchInt = 1
			}
			break
		}
		matchInt = 0
	}

	return matchInt
}

//获取子路由内部的服务来源
func (r *Router) SetServer(server *riderServer.HttpServer) {
	r.server = server
}

//将服务注入子路由
func (r *Router) GetServer() *riderServer.HttpServer {
	return r.server
}

//添加中间件(将根路由服务的中间件加进来)
func (r *Router) AppendMiddleware(path string, rootMiddleware ...HandlerFunc) {
	if r.PathMiddleware == nil {
		r.PathMiddleware = make(map[string][]HandlerFunc)
	}
	if r.mux == nil {
		r.mux = new(sync.Mutex)
	}
	r.mux.Lock()
	r.PathMiddleware[path] = append(rootMiddleware, r.PathMiddleware[path]...)

	r.mux.Unlock()
}

//匹配子路由的方法和路径是否有注册过
func (r *Router) MatchMethodAndPath(method string, path string) int {
	return MatchMethodAndPath(method, path, r.SubRouter)
}

//获取路由根路径
func (r *Router) GetRootPath() string {
	return r.rootPath
}

func (r *Routers) RegisterRouter(method string, path string, router *Router) {
	//赋值子路由的根路径
	router.rootPath = path
	router.currentPath = path
	r.mux.Lock()
	r.routers[path] = router
	r.mux.Unlock()

	//如果这个router没有更深的子集，要处理的子路由也就是其本身了
	//app.GET("/justRouter", &riderRouter.Router{Handler:})
	if len(router.SubRouter) == 0 {
		r.doRouter(method, path, router)
	}

	if router.SubRouter != nil {
		var handleMethod string
		for subPath, handlerRoute := range router.SubRouter {
			for _, handleMap := range handlerRoute {
				//routers里面绑定了app的中间处理，传给router（根路由中的子路由入口点）
				handleMap.AppendMiddleware(subPath, r.PathMiddleware[path]...)
				handleMethod = handleMap.Method
				handleMap.rootMethod = method
				switch handleMethod {
				case "ANY":
					router.ANY(subPath, handleMap)
				case "GET":
					router.GET(subPath, handleMap)
				case "POST":
					router.POST(subPath, handleMap)
				case "HEAD":
					router.HEAD(subPath, handleMap)
				case "DELETE":
					router.DELETE(subPath, handleMap)
				case "PUT":
					router.PUT(subPath, handleMap)
				case "PATCH":
					router.PATCH(subPath, handleMap)
				case "OPTIONS":
					router.OPTIONS(subPath, handleMap)
				case "CONNECT":
					router.CONNECT(subPath, handleMap)
				case "TRACE":
					router.TRACE(subPath, handleMap)
				}
			}
		}
	}
}

//如果一个Router的SubRouter为nil，那么就直接处理这个Router
func (r *Routers) doRouter(method string, path string, router *Router) {
	router.InitParams(method, path, path, path)
	router.AppendMiddleware(path, router.Middleware...)
	router.AppendMiddleware(path, r.PathMiddleware[path]...)
	registeredRouters.NewRouter(method, path, router)
}

//路由服务注册内部入口，由rider的http方法引入
//path：由rider入口出传入，根路由的path（子路由的rootPath）
//router：由rider入口传入，子路由的路由
//在进入ANY内部之前，会先走子路由的ANY方法，将子路由的router和path注册进来
func (r *Routers) ANY(path string, router *Router) {
	r.RegisterRouter("ANY", path, router)
}
func (r *Routers) GET(path string, router *Router) {
	r.RegisterRouter("GET", path, router)
}
func (r *Routers) POST(path string, router *Router) {
	r.RegisterRouter("POST", path, router)
}
func (r *Routers) OPTIONS(path string, router *Router) {
	r.RegisterRouter("OPTIONS", path, router)
}
func (r *Routers) CONNECT(path string, router *Router) {
	r.RegisterRouter("CONNECT", path, router)
}
func (r *Routers) HEAD(path string, router *Router) {
	r.RegisterRouter("HEAD", path, router)
}
func (r *Routers) PUT(path string, router *Router) {
	r.RegisterRouter("PUT", path, router)
}
func (r *Routers) PATCH(path string, router *Router) {
	r.RegisterRouter("PATCH", path, router)
}
func (r *Routers) DELETE(path string, router *Router) {
	r.RegisterRouter("DELETE", path, router)
}
func (r *Routers) TRACE(path string, router *Router) {
	r.RegisterRouter("TRACE", path, router)
}



//重制Router内的所有数据
func (r *Router) Reset() *Router {
	route := NewRouter(nil)
	route.Middleware = r.Middleware
	route.Handler = r.Handler
	route.server = r.server
	return route
}

//http方法的子路由绑定使用，或者根路由通过入口最终进入的实现http服务绑定的地方
func (r *Router) ANY(path string, router *Router) {
	r.RegisterRouter("ANY", path, router)
}
func (r *Router) GET(path string, router *Router) {
	r.RegisterRouter("GET", path, router)
}
func (r *Router) POST(path string, router *Router) {
	r.RegisterRouter("POST", path, router)
}
func (r *Router) OPTIONS(path string, router *Router) {
	r.RegisterRouter("OPTIONS", path, router)
}
func (r *Router) CONNECT(path string, router *Router) {
	r.RegisterRouter("CONNECT", path, router)
}
func (r *Router) HEAD(path string, router *Router) {
	r.RegisterRouter("HEAD", path, router)
}
func (r *Router) PUT(path string, router *Router) {
	r.RegisterRouter("PUT", path, router)
}
func (r *Router) PATCH(path string, router *Router) {
	r.RegisterRouter("PATCH", path, router)
}
func (r *Router) DELETE(path string, router *Router) {
	r.RegisterRouter("DELETE", path, router)
}
func (r *Router) TRACE(path string, router *Router) {
	r.RegisterRouter("TRACE", path, router)
}

//创建中间处理函数
//传入的handler要实现了HttpHandler
func MiddleWare(handlers ...HandlerFunc) []HandlerFunc {
	return handlers
}

//子路由对应的各个http方法的总注册入口
//当多个rootPath路由到r时，由于currentPath相同，
//用currentPath注册会将之前的覆盖掉，由于指针引用，其实引用的都是同一个，最后一次的改变会对之前所有注册同一个路由的router产生影响
//所以在path后暂时加个字符串（@#），然后最后注册时回去掉
func (r *Router) RegisterRouter(method string, path string, router *Router) {


		if ok := r.NewSubRouter(method, path, router); ok {
			return
		}
		// ..//../ -> /
		pattern := filepath.Clean(r.rootPath + path)

		//如果根方法（rider注册的http方法）部位ANY,并且子方法不为ANY && 子方法和根方法不同，将不会注册该路由
		r.matchHTTPMethod(method, pattern, router)

		router.InitParams(method, r.rootPath, pattern, path)

		registeredRouters.NewRouter(method, pattern, router)

		if len(router.PathMiddleware[path]) == 0 && router.Handler == nil {
			log.Println("未被使用的路由：", r.rootPath+path)
			return
		}
	//r.server.ServerMux.Handle(pattern, router)
}

//注册子路由
//如果同名路由重复注册返回false，第一次注册返回true
func (r *Router) NewSubRouter(method string, path string, router *Router) bool {
	//如果子路由没有注册这个path，将其存入根路由入口（Routers.routers.SubRouter）的子路由中
	//如果存在了，说明之前已经注册过了
	registerStatus := r.MatchMethodAndPath(method, path)
	if registerStatus == 0 || registerStatus == 1 {
		router.Method = method
		r.mux.Lock()
		if r.SubRouter[path] == nil {
			r.SubRouter[path] = make(handlerRouter)
		}

		r.SubRouter[path][method] = router
		r.mux.Unlock()

		//将根路由的中间处理函数插入子路由的中间处理函数的首位
		//router.AppendMiddleware(path, r.PathMiddleware[r.rootPath]...)
		//把中间件添加到对应路由内部，以防，多个rider实例同时绑定一个router时，重复添加
		router.AppendMiddleware(path, router.Middleware...)
		router.AppendMiddleware(path, r.Middleware...)
		return true
	}
	return false
}

//初始化Router的参数
func (r *Router) InitParams(method string, rootPath string, fullPath string, currentPath string) {
	r.Method = method
	r.rootPath = rootPath
	//fullPath：完整的路径
	r.FullPath = fullPath
	//currentPath：调用方法传入的路径，相当于fullPath除去rootPath的后半段
	r.currentPath = currentPath
}

//根路由的http方法和子路由的http方法校对
//如果根方法（rider注册的http方法）部位ANY,并且子方法不为ANY && 子方法和根方法不同，将不会注册该路由
func (r *Router) matchHTTPMethod(method string, pattern string, router *Router) {
	if router.rootMethod != "ANY" {
		if method != "ANY" && method != router.rootMethod {
			log.Fatalln("根路由设置" + router.rootMethod + "方法后无法设置子路由的" + method + "；已忽略" + method + "请求的路由" + pattern + "")
			//continue walkIn
		} else if method == "ANY" {
			method = router.rootMethod
		}
	}
}

func (r *Router) registerRouter(routeMethod string, path string, router *Router) {

}