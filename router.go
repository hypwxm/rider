//package router
//路由注册模块，注册http请求（Get,Post,Delete,Option...via）
package rider

import (
	"net/http"
	"sync"
	"path/filepath"
	"strings"
	"container/list"
	"regexp"
	"rider/logger"
	"time"
	"runtime/debug"
)

//跟路由和子路由需要实现的方法
type BaseRouter interface {
	USE(handlers ...HandlerFunc) //添加中间件
	ANY(path string, router *Router) *Router
	POST(path string, router *Router)
	GetServer() *HttpServer
	SetServer() *HttpServer
	GET(path string, router *Router)
	HEAD(path string, router *Router)
	OPTIONS(path string, router *Router)
	PUT(path string, router *Router)
	PATCH(path string, router *Router)
	DELETE(path string, router *Router)
	//HiJack(path string, router *Router)
	//WebSocket(path string, router *Router)
	Any(path string, router *Router)
	registerHandler(routeMethod string, path string, router *Router)
	registerRouter(method string, path string, router *Router)
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
func (r *Router) matchMethodAndPath(method string, path string) int {
	return matchMethodAndPath(method, path, r.subRouter)
}


//给RegisteredRouters添加新的的处理路由
func (r *Router) registerRouter(method string, pattern string, router *Router) {
	//在router.NewSubRouter()中已经验证了路由的唯一性
	r.mux.Lock()
	r.NewPath(pattern)
	r.subRouter[pattern][method] = router
	r.mux.Unlock()
	router.SetServer(r.GetServer())

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

//判断注册的路由里面和请求的路由是否匹配，包括路由参数匹配
func (r *Router) getByPath(path string, request *Request) handlerRouter {
	if filepath.Clean(path) != "/" && strings.LastIndex(path, "/") == len(path) - 1 {
		//path == "/a/b/c/" 去除最后的"/"在进行比较
		path = path[:len(path) - 1]
	}
	walk:
	for k, v := range r.subRouter {
		if path == k {
			return v
		} else {
			//拆解定义的路由
			params := strings.Split(k, "/")
			//拆解请求的路由
			pathParams := strings.Split(path, "/")

			//len(pathParams) >= len(params)

			//判断正则匹配
			//cleanPath := strings.Replace(k, "/", "\\/", -1)
			//fmt.Println(cleanPath)
			reg, err := regexp.Compile("^" + k + "$")
			if err != nil {
				panic(err)
			}

			//定义的路径转换成正则后满足请求的路径
			if reg.MatchString(path) {
				matchParams := reg.FindAllSubmatch([]byte(path), -1)   //[[00 00 00] [11 11 11]]
				if len(matchParams) == 1 {
					matchParams2 := matchParams[0]  //[00 00 00]
					//索引为0的值是路径本身，不需要, 取matchParams2[1:]
					if len(matchParams2) > 1 {
						for _, pathParamsValue := range matchParams2[1:] {
							request.pathParams = append(request.pathParams, string(pathParamsValue))
						}
					}
				}
				return v
			}

			//比较拆解的两个切片的长度，如果长度不一样，肯定就不匹配
			if len(params) != len(pathParams) {
				continue
			}
			//如果长度相等了，就按顺序匹配每一个字段
			for pIndex, param := range params {
				//  ":"开头的说明是路由参数，
				if !(strings.HasPrefix(param, ":")) {
					if pathParams[pIndex] != param {
						request.params = make(map[string]string)
						continue walk
					} else if pIndex == len(params)-1 {
						return v
					}
				} else {
					paramName := strings.TrimPrefix(param, ":")
					request.params[paramName] = pathParams[pIndex]
					if pIndex == len(params)-1 {
						return v
					}
				}
			}
		}
	}
	return nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	routerStart := time.Now()
	context := newContext(w, req, r.GetServer())
	defer func() {
		statusCode := context.Status()
		if err, ok := recover().(error); ok {
			r.server.logger.ERROR(err.Error())
			r.server.logger.PANIC(string(debug.Stack()))
			statusCode = http.StatusInternalServerError
			HttpError(context, err.Error(), http.StatusInternalServerError)
		}
		//记录http请求耗时
		routerEnd := time.Now()
		routerDuration := routerEnd.Sub(routerStart)
		logger.HttpLogger(r.GetServer().logger, context.Method(), context.Path(), statusCode, routerDuration, context.RequestID())
		releaseContext(context)
	}()

	if handler := r.getByPath(req.URL.Path, context.request); handler != nil {
		handler.riderServeHTTP(context)
	} else {
		HttpError(context, "does not match the request path", 404)
	}
}

func (h handlerRouter) riderServeHTTP(context *Context) {
	req := context.request
	w := context.response
	//请求的路径
	reqPath := req.Path()
	//请求的http方法
	reqMethod := req.Method()

	//处理options
	if strings.ToUpper(reqMethod) == http.MethodOptions {
		allows := h.allow(reqPath)
		w.SetHeader("allow", allows)
		return
	}

	var (
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
		}
	}

	if anyRoute != nil {
		anyRoute.doHandlers(context)
		return
	}

	//当没有任何可匹配的请求响应时，用全局的server（所有的router保存的都是全局的server）来响应404
	HttpError(context, "该路由不支持该HTTP方法", 400)
}

//请求方式匹配后将执行该路由注册的handlers和middleware
func (r *Router) doHandlers(context *Context) {
	//未注册任何中间件和处理方式，直接返回404
	if len(r.Middleware) == 0 && r.Handler == nil {
		HttpError(context, "registered but been not used", 404)
		return
	}
	//执行中间件的链表处理，必须调用context.Next()才能执行下一个处理函数
	if r.handlerList != nil && r.handlerList.Len() > 0 {
		context.handlerList = r.handlerList
		context.startHandleList()
	}
}

//调用使用添加中间件(将app服务的中间件加进来)
func (r *Router) USE(middleware ...HandlerFunc) {
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
func matchMethodAndPath(method string, path string, registeredRouters RegisteredRouters) int {
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
	r.registerHandler("ANY", path, router)
}
func (r *Router) GET(path string, router *Router) {
	r.registerHandler("GET", path, router)
}
func (r *Router) POST(path string, router *Router) {
	r.registerHandler("POST", path, router)
}
func (r *Router) OPTIONS(path string, router *Router) {
	r.registerHandler("OPTIONS", path, router)
}
func (r *Router) CONNECT(path string, router *Router) {
	r.registerHandler("CONNECT", path, router)
}
func (r *Router) HEAD(path string, router *Router) {
	r.registerHandler("HEAD", path, router)
}
func (r *Router) PUT(path string, router *Router) {
	r.registerHandler("PUT", path, router)
}
func (r *Router) PATCH(path string, router *Router) {
	r.registerHandler("PATCH", path, router)
}
func (r *Router) DELETE(path string, router *Router) {
	r.registerHandler("DELETE", path, router)
}
func (r *Router) TRACE(path string, router *Router) {
	r.registerHandler("TRACE", path, router)
}

//创建中间处理函数
//传入的handler要实现了HttpHandler
func MiddleWare(handlers ...HandlerFunc) []HandlerFunc {
	return handlers
}

//子路由对应的各个http方法的总注册入口
func (r *Router) registerHandler(method string, path string, router *Router) {
	//赋予router的上级方法rootMethod(不代表整个路由的根，指代"router"的调用者"r"的http方法)
	router.rootMethod = method
	//子路由的根路径为当前父级路径加上子路由当前的rootPath，组成新的rootPath(fullPath为rootPath + 注册时自己的路由)
	router.rootPath = filepath.Join(path, router.rootPath)

	if r.isRoot {
		//如果是根路由，代表可以将之前注册的http方法和路由注册进registeredRouters

		if router.subRouter == nil {
			//如果是直接注册的路由，没有引用下级路由
			registerHandler(r, path, router)
		} else {
			//如果路由存在子集路由
			for subPath, handlerRoute := range router.subRouter {
				for _, handleMap := range handlerRoute {
					handleMap.rootMethod = method
					// ..//../ -> /
					pattern := filepath.Join(path, subPath)
					registerHandler(r, pattern, handleMap)
				}
			}

			//注册已经完成，所有的子路由都已经注册在根路由上，子路由的子路由就不再需要了。在这里将其删除释放，减少路由的无用引用
			router.subRouter = nil
		}

	} else {
		if router.subRouter == nil {
			r.newsubRouter(method, path, router)
		} else {
			//赋值子路由的根路径
			//深层次的路由，rootPath会一直被修改，最终达到根路由变成从根路由开始到它自己的上级  相当于  /a/b/c/d的/a/b/c
			router.rootPath = filepath.Join(path, router.rootPath)
			for subPath, handlerRoute := range router.subRouter {
				for _, handleMap := range handlerRoute {
					handleMap.rootMethod = method
					// ..//../ -> /
					pattern := filepath.Join(path, subPath)

					r.newsubRouter(method, pattern, handleMap)
				}
			}

		}
	}
}

func registerHandler(r *Router, pattern string, router *Router) {
	// ..//../ -> /
	pattern = filepath.Clean(pattern)

	method := r.newsubRouter(router.rootMethod, pattern, router)

	r.registerRouter(method, pattern, router)

	if len(router.Middleware) == 0 && router.Handler == nil {
		r.server.logger.WARNING("未被使用的路由：", pattern)
	}
}

//注册子路由
//如果同名路由重复注册返回false，第一次注册返回true
func (r *Router) newsubRouter(method string, path string, router *Router) string {

	if router.Method != "" {
		method = matchHTTPMethod(method, path, router)
	}

	//如果子路由没有注册这个path，将其存入根路由入口（Routers.routers.subRouter）的子路由中
	//如果存在了，说明之前已经注册过了
	registerStatus := r.matchMethodAndPath(method, path)
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
