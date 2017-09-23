//package router
//路由注册模块，注册http请求（Get,Post,Delete,Option...via）
package rider

import (
	"net/http"
	"sync"
	"path/filepath"
	"log"
	"strings"
	"container/list"
	"runtime/debug"
)

//跟路由和子路由需要实现的方法
type BaseRouter interface {
	ANY(path string, router *Router) *Router
	POST(path string, router *Router)
	GetServer() *HttpServer
	SetServer() *HttpServer
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
	RegisterHandler(routeMethod string, path string, router *Router)
	RegisterRouter(method string, path string, router *Router)
	MiddleWare(handlers ...HandlerFunc) []HandlerFunc
}

//子路由配置
type Router struct {
	subRouter   RegisteredRouters //子路由还可以继续定义子路由
	isRoot      bool
	mux         *sync.Mutex
	Handler     HandlerFunc
	Middleware  []HandlerFunc
	Method      string      //定义的请求的方法
	fullPath    string      //定义的完整路径
	server      *HttpServer //全局的http服务入口，有rider传入
	rootPath    string      //根路由   /home/banner中的/home
	rootMethod  string      //跟路由上定义的方法(和Method要对应)；要求跟路由为ANY时，子路由请求方法无限制，跟路由不为ANY子路由只能是ANY和同名http方法
	handlerList *list.List  //中间件加处理函数链表（会存放在context中)
}

type handlerRouter map[string]*Router

//  map["/home"]["GET"]*Router
type RegisteredRouters map[string]handlerRouter

func (h handlerRouter) getServer() *HttpServer {
	for _, v := range h {
		return v.server
	}
	return nil
}

func (r *Router) NewPath(path string) {
	if r.subRouter[path] == nil {
		r.mux.Lock()
		r.subRouter[path] = make(handlerRouter)
		r.mux.Unlock()
	}
}

//匹配路由是否已经注册过
func (r *Router) MatchMethodAndPath(method string, path string) int {
	return MatchMethodAndPath(method, path, r.subRouter)
}

//给RegisteredRouters添加新的的处理路由
func (r *Router) RegisterRouter(method string, pattern string, router *Router) {
	//在router.NewSubRouter()中已经验证了路由的唯一性
	r.mux.Lock()
	r.NewPath(pattern)
	r.subRouter[pattern][method] = router
	r.mux.Unlock()

	//形成中间件加主处理函数链表
	router.handlerList = list.New()

	if len(router.Middleware) > 0 {
		for _, middleware := range router.Middleware {
			router.handlerList.PushBack(middleware)
		}
	}
	if router.Handler != nil {
		router.handlerList.PushBack(router.Handler)
	}
}

//子路由初始化方法
func NewRouter() *Router {
	return &Router{
		subRouter: make(RegisteredRouters),
		mux:       new(sync.Mutex),
	}
}

func (r *Router) getByPath(path string, request *Request) handlerRouter {
	var paramsRouter handlerRouter
	walk:
	for k, v := range r.subRouter {
		if path == k {
			return v
		} else {
			params := strings.Split(k, "/")
			pathParams := strings.Split(path, "/")
			if len(params) != len(pathParams) {
				continue
			}
			for pIndex, param := range params {
				if !(strings.HasPrefix(param, ":")) {
					if pathParams[pIndex] != param {
						continue walk
					}
				}
				paramName := strings.TrimPrefix(param, ":")
				request.params[paramName] = pathParams[pIndex]
			}
			paramsRouter = v
		}
	}
	if paramsRouter != nil {
		return paramsRouter
	}
	return nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	defer func() {
		if err, ok := recover().(error); ok {
			http.Error(w, "系统错误", http.StatusInternalServerError)
			log.Println(err)
			log.Println(string(debug.Stack()))
		}
	}()
	response := basePool.response.Get().(*Response)
	request := basePool.request.Get().(*Request)
	context := basePool.context.Get().(*Context)
	response.reset(w)
	request.reset(req)
	context.reset(response, request)

	if handler := r.getByPath(req.URL.Path, request); handler != nil {
		handler.RiderServeHTTP(context)
	} else {
		response.Send("not found")
	}

	response.release()
	basePool.response.Put(response)
	request.release()
	basePool.request.Put(request)
	context.release()
	basePool.context.Put(context)
}

func (h handlerRouter) RiderServeHTTP(context *Context) {
	req := context.Request
	w := context.Response
	//请求的路径
	reqPath := req.Path()
	//请求的http方法
	reqMethod := req.Method()
	if strings.ToUpper(reqMethod) == http.MethodOptions {
		allows := h.allow(reqPath)
		w.SetHeader("allow", allows)
		return
	}

	var (
		finalRouter *Router
		anyRoute    *Router
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
				router.doHandlers(context)
				return
			}
			finalRouter = router
		}
	}

	if anyRoute != nil {
		anyRoute.doHandlers(context)
		return
	}

	//当没有任何可匹配的请求响应时，用全局的server（所有的router保存的都是全局的server）来响应404
	finalRouter.server.Error(w.writer, "not found", http.StatusNotFound)
}

//请求方式匹配后将执行该路由注册的handlers和middleware
func (r *Router) doHandlers(context *Context) {
	//未注册任何中间件和处理方式，直接返回404
	if len(r.Middleware) == 0 && r.Handler == nil {
		r.server.Error(context.Response.writer, "not found", http.StatusNotFound)
		return
	}
	//执行中间件的链表处理，必须调用context.Next()才能执行下一个处理函数
	if r.handlerList != nil && r.handlerList.Len() > 0 {
		context.handlerList = r.handlerList
		context.StartHandleList()
	}
}

//调用使用添加中间件(将app服务的中间件加进来)
func (r *Router) AddMiddleware(middleware ...HandlerFunc) {
	r.Middleware = append(r.Middleware, middleware...)
}

//当请求为options时，返回该路径所支持的请求方式
func (h handlerRouter) allow(path string) string {
	var allows = ""
	for _method, _ := range h {
		if _method == "OPTIONS" {
			continue
		}
		if allows == "" {
			allows += _method
		} else {
			allows += "," + _method
		}
	}
	return allows
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
func (r *Router) SetServer(server *HttpServer) {
	r.server = server
}

//将服务注入子路由
func (r *Router) GetServer() *HttpServer {
	return r.server
}

//添加中间件(将根路由服务的中间件加进来)
func (r *Router) FrontMiddleware(rootMiddleware ...HandlerFunc) {
	r.Middleware = append(rootMiddleware, r.Middleware...)
}

//获取路由根路径
func (r *Router) GetRootPath() string {
	return r.rootPath
}

//获取完整路径
func (r *Router) Path() string {
	return r.fullPath
}

//http方法的子路由绑定使用，或者根路由通过入口最终进入的实现http服务绑定的地方
func (r *Router) ANY(path string, router *Router) {
	r.RegisterHandler("ANY", path, router)
}
func (r *Router) GET(path string, router *Router) {
	r.RegisterHandler("GET", path, router)
}
func (r *Router) POST(path string, router *Router) {
	r.RegisterHandler("POST", path, router)
}
func (r *Router) OPTIONS(path string, router *Router) {
	r.RegisterHandler("OPTIONS", path, router)
}
func (r *Router) CONNECT(path string, router *Router) {
	r.RegisterHandler("CONNECT", path, router)
}
func (r *Router) HEAD(path string, router *Router) {
	r.RegisterHandler("HEAD", path, router)
}
func (r *Router) PUT(path string, router *Router) {
	r.RegisterHandler("PUT", path, router)
}
func (r *Router) PATCH(path string, router *Router) {
	r.RegisterHandler("PATCH", path, router)
}
func (r *Router) DELETE(path string, router *Router) {
	r.RegisterHandler("DELETE", path, router)
}
func (r *Router) TRACE(path string, router *Router) {
	r.RegisterHandler("TRACE", path, router)
}

//创建中间处理函数
//传入的handler要实现了HttpHandler
func MiddleWare(handlers ...HandlerFunc) []HandlerFunc {
	return handlers
}

//子路由对应的各个http方法的总注册入口
func (r *Router) RegisterHandler(method string, path string, router *Router) {
	//赋予router的上级方法rootMethod(不代表整个路由的根，指代"router"的调用者"r"的http方法)
	router.rootMethod = method
	//子路由的根路径为当前父级路径加上子路由当前的rootPath，组成新的rootPath(fullPath为rootPath + 注册时自己的路由)
	router.rootPath = filepath.Join(path, router.rootPath)

	if r.isRoot {
		//如果是根路由，代表可以将之前注册的http方法和路由注册进registeredRouters

		if router.subRouter == nil {
			//如果是直接注册的路由，没有引用下级路由
			RegisterHandler(r, path, router)
		} else {
			//如果路由存在子集路由
			for subPath, handlerRoute := range router.subRouter {
				for _, handleMap := range handlerRoute {
					handleMap.rootMethod = method
					// ..//../ -> /
					pattern := filepath.Join(path, subPath)
					RegisterHandler(r, pattern, handleMap)
				}
			}

			//注册已经完成，所有的子路由都已经注册在根路由上，子路由的子路由就不再需要了。在这里将其删除释放，减少路由的无用引用
			router.subRouter = nil
		}

	} else {
		if router.subRouter == nil {
			r.NewsubRouter(method, path, router)
		} else {
			//赋值子路由的根路径
			//深层次的路由，rootPath会一直被修改，最终达到根路由变成从根路由开始到它自己的上级  相当于  /a/b/c/d的/a/b/c
			router.rootPath = filepath.Join(path, router.rootPath)
			for subPath, handlerRoute := range router.subRouter {
				for _, handleMap := range handlerRoute {
					handleMap.rootMethod = method
					// ..//../ -> /
					pattern := filepath.Join(path, subPath)

					r.NewsubRouter(method, pattern, handleMap)
				}
			}

		}
	}
}

func RegisterHandler(r *Router, pattern string, router *Router) {
	// ..//../ -> /
	pattern = filepath.Clean(pattern)

	method := r.NewsubRouter(router.rootMethod, pattern, router)

	r.RegisterRouter(method, pattern, router)

	if len(router.Middleware) == 0 && router.Handler == nil {
		log.Println("未被使用的路由：", pattern)
	}
}

//注册子路由
//如果同名路由重复注册返回false，第一次注册返回true
func (r *Router) NewsubRouter(method string, path string, router *Router) string {

	if router.Method != "" {
		method = matchHTTPMethod(method, path, router)
	}

	//如果子路由没有注册这个path，将其存入根路由入口（Routers.routers.subRouter）的子路由中
	//如果存在了，说明之前已经注册过了
	registerStatus := r.MatchMethodAndPath(method, path)
	if registerStatus == 0 || registerStatus == 1 {
		router.fullPath = path
		router.Method = method
		r.mux.Lock()
		if r.subRouter[path] == nil {
			r.subRouter[path] = make(handlerRouter)
		}

		r.subRouter[path][method] = router
		r.mux.Unlock()

		//将根路由的中间处理函数插入子路由的中间处理函数的首位
		router.FrontMiddleware(r.Middleware...)
	} else if registerStatus == 2 {
		panic(method + " " + path + " 已注册，请勿重复注册同一个http方法和请求路径")
	}
	return method
}

//根路由的http方法和子路由的http方法校对
//如果根方法（rider注册的http方法）部位ANY,并且子方法不为ANY && 子方法和根方法不同，将不会注册该路由
func matchHTTPMethod(parentMethod string, pattern string, router *Router) string {
	if parentMethod != "ANY" {
		if router.Method != "ANY" && router.Method != parentMethod {
			panic("the parent http method " + parentMethod + " can not set child http method width " + router.Method + "；router " + pattern)
		} else {
			//如果子路由的方法为ANY,就直接使用根路由的方法作为HTTP方法
			return parentMethod
		}
	} else {
		return router.Method
	}
}
