package rider

import (
	"rider/riderRouter"
	"rider/riderServer"
	"net/http"
)

//http服务的入口，用户初始化和缓存服务的一些信息
type rider struct {
	server *riderServer.HttpServer //注册服务用的serveMu，全局统一
	Middleware []riderRouter.HandlerFunc //全局的中间件，惠安循序
	routers *riderRouter.Routers
	Func riderRouter.HandlerFunc
	Method string
}

//初始化服务入口组建
func New() *rider {
	server := &riderServer.HttpServer{ServerMux: http.NewServeMux()}
	return &rider{
		server: server,
		routers: riderRouter.NewRouters(),
	}
}

func (r *rider) NewRouter() *riderRouter.Router {
	_router := riderRouter.NewRouter(nil)
	_router.SetServer(r.server)
	return _router
}

func (r *rider) Run(port string) {
	r.routers.Run()
}

//提供端口监听服务，监听rider里面的serveMux,调用http自带的服务启用方法
func (r *rider) Listen(port string) {
	r.Run(port)
	http.ListenAndServe(port, r.server.ServerMux)
}

//http请求的方法的入口（ANY, GET, POST...VIA）
//path：一个跟路径，函数内部根据这个根路径创建一个根路由routers，用来管理router子路由
//router：这个根路径对应的子路由入口。
func (r *rider) registerRiderRouter(method string, path string, router interface{}) {
	if handleFunc, ok := router.(func(w http.ResponseWriter, r *http.Request)); ok {
		registeredRouters := r.routers.GetRegisteredRouters()
		_router := riderRouter.NewRouter(handleFunc)
		_router.InitParams(method, path, path, path)
		registeredRouters.NewRouter(method, path, _router)
		return
	}

	if router, ok := router.(*riderRouter.Router); ok {
		//将app的中间处理函数传给routers根路由(后面再由routers传给其各个子路由)
		r.routers.AppendMiddleware(path, r.Middleware...)
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
}

func (r *rider) ANY(path string, router interface{}) {
	r.registerRiderRouter("ANY", path, router)
}

func (r *rider) GET(path string, router interface{}) {
	r.registerRiderRouter("GET", path, router)
}

func (r *rider) POST(path string, router interface{}) {
	r.registerRiderRouter("POST", path, router)
}

func (r *rider) HEAD(path string, router interface{}) {
	r.registerRiderRouter("HEAD", path, router)
}

func (r *rider) OPTIONS(path string, router interface{}) {
	r.registerRiderRouter("OPTIONS", path, router)
}

func (r *rider) DELETE(path string, router interface{}) {
	r.registerRiderRouter("DELETE", path, router)
}

func (r *rider) PUT(path string, router interface{}) {
	r.registerRiderRouter("PUT", path, router)
}

func (r *rider) PATCH(path string, router interface{}) {
	r.registerRiderRouter("PATCH", path, router)
}

func (r *rider) TRACE(path string, router interface{}) {
	r.registerRiderRouter("TRACE", path, router)
}

func (r *rider) CONNECT(path string, router interface{}) {
	r.registerRiderRouter("CONNECT", path, router)
}


//返回服务内部的http服务入口
func (r *rider) GetServer() *riderServer.HttpServer {
	return r.server
}

//为app服务添加中间处理
func (r *rider) AddMiddleware(handlers ...riderRouter.HandlerFunc) {
	r.Middleware = append(r.Middleware, handlers...)
}

func (r *rider) Handle(pattern string, handle http.Handler) {
	r.server.ServerMux.Handle(pattern, handle)
}

func (r *rider) HandleFunc(pattern string, handleFunc riderRouter.HandlerFunc) {
	router := riderRouter.NewRouter(handleFunc)
	r.routers.RegisterRouter("ANY", pattern, router)
}

//直接在跟路由处理http请求
func (r *rider) DoFunc(w http.ResponseWriter, req *http.Request) {
	if r.Func != nil {
		r.Func(w, req)
	}
}


func (r *rider) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.DoFunc(w, req)
}