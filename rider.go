/*
*	路由为内部主要模块，不提供对用实现接口，请遵照一下要求：
*	1: 路由同名者，同HTTP路由重复注册，会在注册时panic；
*	2: 同名路由存在ANY以外的方法时，会选择客户端请求的方法作为响应，当没有其他方法，会响应ANY方法
*	3: 多层级子路由注册时，必须保持子路由的HTTP方法和父级路由的关系为（1:相同； 2:其中一个或两个为ANY）否者panic
*	4: 中间件注册统一为路由实例的AddMiddleware方法（中间分为根路由中间件，子路由中间件，和路由内部中间件），详情参阅example的middleware模块
	5: 调用下一个中间件为context.Next()，不调用将不会继续往下执行。
*	6: 服务启动入口为rider.Listen(port:":8000")；默认端口":8000"；
*	7: 所有的路由都是维护在一个初始化的变量内部registeredRouters;
*	8: 支持无限极子路由，会一层一层的检测，当前层是否会有1，2情况的发生。如有1，2的情况发生，请注意错误，不会定位到完整路由。
*/

package rider

import (
	"net/http"
	"log"
	"time"
)

const (
	addr           = ":8000"
	readTimeout    = 30 * time.Second
	writerTimeout  = 5 * time.Second
	maxHeaderBytes = 1 << 20 //1MB
	defaultMultipartBodySze = 32 << 20
	ENV_Production = "production"
	ENV_Development = "development"
	ENV_Debug = "debug"
)

//默认的执行的系统的环境，生产者环境
var GlobalENV = ENV_Production

type baseRider interface {
	Listen(port string) //:5000  服务启动入口
	registerRiderRouter(method string, path string, router *Router)
	ANY(path string, router *Router)
	GET(path string, router *Router)
	POST(path string, router *Router)
	HEAD(path string, router *Router)
	OPTIONS(path string, router *Router)
	DELETE(path string, router *Router)
	PUT(path string, router *Router)
	PATCH(path string, router *Router)
	TRACE(path string, router *Router)
	CONNECT(path string, router *Router)
	GetServer() *HttpServer   //获取
	AddMiddleware(handlers ...HandlerFunc)
	Error(errorHandle func(c *Context, err string, code int))
}

//http服务的入口，用户初始化和缓存服务的一些信息
type rider struct {
	server   *HttpServer //注册服务用的serveMu，全局统一
	routers  *Router
}

//设置环境
func SetEnvMode(mode string) {
	if mode == "production" {
		GlobalENV = ENV_Production
	} else if mode == "development" {
		GlobalENV = ENV_Development
	} else if mode == "debug" {
		GlobalENV = ENV_Debug
	}
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

//提供端口监听服务，监听rider里面的serveMux,调用http自带的服务启用方法
func (r *rider) Listen(port string) {

	if port == "" {
		port = addr
	}
	server := &http.Server{
		Addr:           port,
		Handler:        r.routers,
		ReadTimeout:    readTimeout,
		WriteTimeout:   0,
		MaxHeaderBytes: maxHeaderBytes,
	}
	err := server.ListenAndServe()
	//err := http.ListenAndServe(port, r.server.ServerMux)
	if err != nil {
		log.Fatalln("[error] ", err)
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


//重写错误处理
func (r *rider) Error(errorHandle func(c *Context, err string, code int)) {
	ErrorHandle = errorHandle
}