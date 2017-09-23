package rider

import (
	"net/http"
	"log"
	"time"
)

const (
	addr           = ":8000"
	readTimeout    = 30 * time.Second
	writerTimeout  = 30 * time.Second
	maxHeaderBytes = 1 << 20 //1MB
)

type baseRider interface {
	Listen(port string) //:5000
	Run()               //:5000
}

//http服务的入口，用户初始化和缓存服务的一些信息
type rider struct {
	server     *HttpServer   //注册服务用的serveMu，全局统一
	routers    *Router
}

//初始化服务入口组建
func New() *rider {
	server := &HttpServer{ServerMux: http.NewServeMux()}
	return &rider{
		server:  server,
		routers: NewRootRouter(),
	}
}

func NewRootRouter() *Router {
	_router := NewRouter()
	_router.isRoot = true
	_router.fullPath = "/"
	_router.rootPath = "/"
	_router.Method = "ANY"
	return _router
}

func (r *rider) Run() {
	r.routers.Run()
}

//提供端口监听服务，监听rider里面的serveMux,调用http自带的服务启用方法
func (r *rider) Listen(port string) {

	if port == "" {
		port = addr
	}
	r.Run()
	server := &http.Server{
		Addr:           port,
		Handler:        r.server.ServerMux,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writerTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}
	err := server.ListenAndServe()
	//err := http.ListenAndServe(port, r.server.ServerMux)
	if err != nil {
		log.Fatalln(err)
	}
}

//http请求的方法的入口（ANY, GET, POST...VIA）
//path：一个跟路径，函数内部根据这个根路径创建一个根路由routers，用来管理router子路由
//router：这个根路径对应的子路由入口。
func (r *rider) registerRiderRouter(method string, path string, router *Router) {
	//将app的中间处理函数传给routers根路由(后面再由routers传给其各个子路由)
	//r.routers.FrontMiddleware(r.Middleware...)
	//将服务注入这个组建根路由，确保路由是创立在这个服务上
	r.routers.SetServer(r.server)
	//将服务注入这个组建子路由，确保路由是创立在这个服务上
	router.SetServer(r.server)
	//服务注册入口内部的入口
	switch method {
	case "ANY":
		r.routers.ANY(path, router)
	case "GET":
		r.routers.GET(path, router)
	case "POST":
		r.routers.POST(path, router)
	case "HEAD":
		r.routers.HEAD(path, router)
	case "DELETE":
		r.routers.DELETE(path, router)
	case "PUT":
		r.routers.PUT(path, router)
	case "PATCH":
		r.routers.PATCH(path, router)
	case "OPTIONS":
		r.routers.OPTIONS(path, router)
	case "CONNECT":
		r.routers.CONNECT(path, router)
	case "TRACE":
		r.routers.TRACE(path, router)
	}
}

func (r *rider) ANY(path string, router *Router) {
	r.registerRiderRouter("ANY", path, router)
}

func (r *rider) GET(path string, router *Router) {
	r.registerRiderRouter("GET", path, router)
}

func (r *rider) POST(path string, router *Router) {
	r.registerRiderRouter("POST", path, router)
}

func (r *rider) HEAD(path string, router *Router) {
	r.registerRiderRouter("HEAD", path, router)
}

func (r *rider) OPTIONS(path string, router *Router) {
	r.registerRiderRouter("OPTIONS", path, router)
}

func (r *rider) DELETE(path string, router *Router) {
	r.registerRiderRouter("DELETE", path, router)
}

func (r *rider) PUT(path string, router *Router) {
	r.registerRiderRouter("PUT", path, router)
}

func (r *rider) PATCH(path string, router *Router) {
	r.registerRiderRouter("PATCH", path, router)
}

func (r *rider) TRACE(path string, router *Router) {
	r.registerRiderRouter("TRACE", path, router)
}

func (r *rider) CONNECT(path string, router *Router) {
	r.registerRiderRouter("CONNECT", path, router)
}

//返回服务内部的http服务入口
func (r *rider) GetServer() *HttpServer {
	return r.server
}

//为app服务添加中间处理
func (r *rider) AddMiddleware(handlers ...HandlerFunc) {
	r.routers.Middleware = append(r.routers.Middleware, handlers...)
}
