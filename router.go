//package router
//路由注册模块，注册http请求（Get,Post,Delete,Option...via）
package rider

import (
	"container/list"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/hypwxm/rider/logger"
)

var sameRouterError error = errors.New("duplicate route registration")
var emptyRouterError error = errors.New("empty router")

type routerError struct {
	path   string
	method string
	error  error
}

func (re routerError) string() string {
	return re.method + " " + re.path + " " + re.error.Error()
}

//跟路由和子路由需要实现的方法
type BaseRouter interface {
	USE(handlers ...HandlerFunc) //添加中间件
	ANY(path string, handlers ...HandlerFunc)
	POST(path string, handlers ...HandlerFunc)
	GetServer() *HttpServer
	SetServer(*HttpServer)
	GET(path string, handlers ...HandlerFunc)
	HEAD(path string, handlers ...HandlerFunc)
	OPTIONS(path string, handlers ...HandlerFunc)
	PUT(path string, handlers ...HandlerFunc)
	PATCH(path string, handlers ...HandlerFunc)
	DELETE(path string, handlers ...HandlerFunc)
	Kid(path string, router *Router)
	registerHandler(routeMethod string, path string, router *Router)
	registerRouter(method string, path string, router *Router)
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
	handlerList *list.List  //中间件加处理函数链表（会存放在context中)
}

type handlerRouter map[string]*Router

//  map["/home"]["GET"]*Router
type RegisteredRouters map[string]handlerRouter

//子路由初始化方法
func NewRouter() *Router {
	return &Router{
		subRouter: make(RegisteredRouters),
		mux:       new(sync.Mutex),
	}
}

//判断注册的路由里面和请求的路由是否匹配，包括路由参数匹配（路由不区分大小写）
/**
匹配顺序（正则匹配和params请不要混合使用）
1:完整匹配
2:params匹配
3:正则匹配
**/
func (r *Router) getByPath(path string, request *Request) handlerRouter {
	// windows下的路径\转成/
	path = strings.Replace(path, "\\", "/", -1)
	if filepath.Clean(path) != "/" && strings.LastIndex(path, "/") == len(path)-1 {
		//path == "/a/b/c/" 去除最后的"/"在进行比较
		path = path[:len(path)-1]
	}
	// 先进行绝对匹配 （可以加快绝对匹配速度，/a/b和/a/:id可以共存，除非id==b，否则不会冲突）
	for k, v := range r.subRouter {
		if strings.ToLower(path) == strings.ToLower(k) {
			return v
		}
	}
	//walk:
	// 在进行非绝对匹配，（定义了路径参数的路由）
	for k, v := range r.subRouter {

		//拆解定义的路由
		params := strings.Split(k, "/")
		//拆解请求的路由
		pathParams := strings.Split(path, "/")

		//len(pathParams) >= len(params)

		//比较拆解的两个切片的长度，如果长度不一样，肯定就不匹配
		if len(params) == len(pathParams) {

			//如果长度相等了，就按顺序匹配每一个字段
			for pIndex, param := range params {
				//  ":"开头的说明是路由参数，
				if !(strings.HasPrefix(param, ":")) {
					if strings.ToLower(pathParams[pIndex]) != strings.ToLower(param) {
						request.params = make(map[string]string)
						break
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

		//判断正则匹配
		//cleanPath := strings.Replace(k, "/", "\\/", -1)
		//fmt.Println(cleanPath)
		reg, err := regexp.Compile("^" + k + "$")
		if err != nil {
			panic(err)
		}

		//定义的路径转换成正则后满足请求的路径
		if reg.MatchString(path) {
			matchParams := reg.FindAllSubmatch([]byte(path), -1) //[[00 00 00] [11 11 11]]
			if len(matchParams) == 1 {
				matchParams2 := matchParams[0] //[00 00 00]
				//索引为0的值是路径本身，不需要, 取matchParams2[1:]
				if len(matchParams2) > 1 {
					for _, pathParamsValue := range matchParams2[1:] {
						request.pathParams = append(request.pathParams, string(pathParamsValue))
					}
				}
			}
			return v
		}
	}
	return nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	routerStart := time.Now()
	ctx := newContext(w, req, r.GetServer())
	defer func() {
		statusCode := ctx.Status()
		if err, ok := recover().(error); ok {
			r.server.logger.ERROR(err.Error())
			r.server.logger.PANIC(string(debug.Stack()))
			statusCode = http.StatusInternalServerError
			HttpError(ctx, err.Error(), http.StatusInternalServerError)
		}
		//记录http请求耗时
		routerEnd := time.Now()
		routerDuration := routerEnd.Sub(routerStart)
		logger.HttpLogger(r.GetServer().logger, ctx.Method(), ctx.Request().RequestURI(), statusCode, routerDuration, ctx.RequestID())
		releaseContext(ctx)
	}()

	if handler := r.getByPath(req.URL.Path, ctx.Request()); handler != nil {
		handler.riderServeHTTP(ctx)
	} else {
		HttpError(ctx, "does not match the request path", 404)
	}
}

func (h handlerRouter) riderServeHTTP(context Context) {
	req := context.Request()
	w := context.Response()
	//请求的路径
	reqPath := req.Path()
	//请求的http方法
	reqMethod := req.Method()

	// 此部分为跨域设置（要在option之前做处理，不然patch，put等请求不能绕过option请求会导致跨域失败）
	httpserver := context.getHttpServer()
	access := httpserver.accessControl(context)
	if access.AccessControlAllowCredentials != "" {
		context.SetHeader("Access-Control-Allow-Credentials", access.AccessControlAllowCredentials)
	}

	if access.AccessControlAllowHeaders != "" {
		context.SetHeader("Access-Control-Allow-Headers", access.AccessControlAllowHeaders)
	}

	if access.AccessControlAllowMethods != "" {
		context.SetHeader("Access-Control-Allow-Methods", access.AccessControlAllowMethods)
	}

	if access.AccessControlAllowOrigin != "" {
		context.SetHeader("Access-Control-Allow-Origin", access.AccessControlAllowOrigin)
	}

	//处理options
	if strings.ToUpper(reqMethod) == http.MethodOptions {
		allows := h.allow(reqPath)
		w.SetHeader("allow", allows)
		return
	}

	var (
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
func (r *Router) doHandlers(context Context) {
	//未注册任何中间件和处理方式，直接返回404
	if len(r.Middleware) == 0 && r.Handler == nil {
		HttpError(context, "registered but been not used", 404)
		return
	}
	//执行中间件的链表处理，必须调用context.Next()才能执行下一个处理函数
	if r.handlerList != nil && r.handlerList.Len() > 0 {
		context.handlerQueue(r.handlerList)
		context.startHandleList()
	}
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

//获取完整路径
func (r *Router) Path() string {
	return r.fullPath
}

//调用使用添加中间件(将app服务的中间件加进来)
func (r *Router) USE(middleware ...HandlerFunc) {
	r.Middleware = append(r.Middleware, middleware...)
}

//http方法的子路由绑定使用，或者根路由通过入口最终进入的实现http服务绑定的地方
func (r *Router) ANY(path string, handlers ...HandlerFunc) {
	r.addRouter("ANY", path, handlers...)
}
func (r *Router) GET(path string, handlers ...HandlerFunc) {
	r.addRouter("GET", path, handlers...)
}
func (r *Router) POST(path string, handlers ...HandlerFunc) {
	r.addRouter("POST", path, handlers...)
}
func (r *Router) OPTIONS(path string, handlers ...HandlerFunc) {
	r.addRouter("OPTIONS", path, handlers...)
}
func (r *Router) CONNECT(path string, handlers ...HandlerFunc) {
	r.addRouter("CONNECT", path, handlers...)
}
func (r *Router) HEAD(path string, handlers ...HandlerFunc) {
	r.addRouter("HEAD", path, handlers...)
}
func (r *Router) PUT(path string, handlers ...HandlerFunc) {
	r.addRouter("PUT", path, handlers...)
}
func (r *Router) PATCH(path string, handlers ...HandlerFunc) {
	r.addRouter("PATCH", path, handlers...)
}
func (r *Router) DELETE(path string, handlers ...HandlerFunc) {
	r.addRouter("DELETE", path, handlers...)
}
func (r *Router) TRACE(path string, handlers ...HandlerFunc) {
	r.addRouter("TRACE", path, handlers...)
}
func (r *Router) Kid(path string, router ...IsRouterHandler) {
	r.addChild(path, router...)
}

//子路由对应的各个http方法的总注册入口
//只有r为根路由app时，才能够注册路由，子路由调用此方法需先拼接路由，待拼接到app位置
func (r *Router) addRouter(method string, path string, handlers ...HandlerFunc) error {
	r.mux.Lock()
	if r.subRouter == nil {
		r.subRouter = make(RegisteredRouters)
	}

	// 去掉路径后的/
	path = filepath.Clean(path)

	if r.subRouter[path] == nil {
		r.subRouter[path] = make(handlerRouter)
	}

	//判断路由注册是否重复
	r.duplicateRoute(path, method)

	route := &Router{}
	hLen := len(handlers)
	switch hLen {
	case 0:
		route.Middleware = r.Middleware
	case 1:
		route.Middleware = r.Middleware
		route.Handler = handlers[0]
	default:
		route.Middleware = append(r.Middleware, handlers[:hLen-1]...)
		route.Handler = handlers[hLen-1]
	}

	route.fullPath = path
	route.Method = method

	//根路由处理开始注册路由
	if r.isRoot {
		//注册路由服务
		route.SetServer(r.server)
		//形成中间件加主处理函数链表
		route.handlerList = list.New()

		if len(route.Middleware) > 0 {
			for _, middleware := range route.Middleware {
				route.handlerList.PushBack(middleware)
			}
		}
		if route.Handler != nil {
			route.handlerList.PushBack(route.Handler)
		}

	}

	r.subRouter[path][method] = route

	r.mux.Unlock()
	return nil
}

func (r *Router) addChild(path string, router ...IsRouterHandler) {

	if len(router) == 0 {
		r.server.logger.FATAL("nothing todo with no handler router")
	}

	var funcMid []HandlerFunc

	var toRouter *Router

	for k, _ := range router {
		if k < len(router)-1 {
			if route, ok := router[k].(HandlerFunc); !ok {
				r.server.logger.FATAL("put kid router at the end when register kid router")
			} else {
				funcMid = append(funcMid, route)
			}
		} else {
			if route, ok := router[k].(*Router); !ok {
				r.server.logger.FATAL("if you want register handler directly do not use kid")
			} else {
				//将服务注入这个组建子路由，确保路由是创立在这个服务上
				route.SetServer(r.server)
				toRouter = route

				if len(funcMid) > 0 {
					toRouter.FrontMiddleware(funcMid[:]...)
				}
			}
		}
	}

	if r.isRoot {
		//如果是根路由，代表可以将之前注册的http方法和路由注册进registeredRouters
		if toRouter.subRouter == nil {
			//如果是直接注册的路由，请使用http方法直接注册
			r.server.logger.FATAL("please use GET, POST via http method register router " + path)
		} else {
			//如果路由存在子集路由
			for subPath, methodMap := range toRouter.subRouter {
				for _method, methodRouter := range methodMap {
					// ..//../ -> /
					fullPath := filepath.Join(path, subPath)
					//深层次的路由，rootPath会一直被修改，最终达到根路由变成从根路由开始到它自己的上级  相当于  /a/b/c/d的/a/b/c
					methodRouter.fullPath = fullPath

					r.addRouter(_method, fullPath, append(funcMid[:], append(methodRouter.Middleware, methodRouter.Handler)...)...)
				}
			}
			//注册已经完成，所有的子路由都已经注册在根路由上，子路由的子路由就不再需要了。在这里将其删除释放，减少路由的无用引用
			toRouter = nil
		}
	} else {
		if toRouter.subRouter == nil {
			log.Println("[ERROR] please use GET, POST via http method register router " + path)
			os.Exit(1)
		} else {
			//赋值子路由的根路径
			for subPath, methodMap := range toRouter.subRouter {
				for _method, methodRouter := range methodMap {
					// ..//../ -> /
					fullPath := filepath.Join(path, subPath)

					r.mux.Lock()
					//如果子路由没有注册这个path，将其存入根路由入口（Routers.routers.subRouter）的子路由中
					//如果存在了，说明之前已经注册过了
					r.duplicateRoute(fullPath, _method)

					methodRouter.fullPath = fullPath

					if r.subRouter[fullPath] == nil {
						r.subRouter[fullPath] = make(handlerRouter)
					}

					//将根路由的中间处理函数插入子路由的中间处理函数的首位
					methodRouter.FrontMiddleware(append(r.Middleware, funcMid[:]...)...)
					r.subRouter[fullPath][_method] = methodRouter
					r.mux.Unlock()
				}
			}

		}
	}
}

//判断路由重复
func (r *Router) duplicateRoute(path string, method string) {
	if methodRouter, ok := r.subRouter[path]; ok {
		if _, ok := methodRouter[method]; ok {
			log.Println("[ERROR] ", routerError{path, method, sameRouterError}.string())
			os.Exit(1)
		}
	}
}
