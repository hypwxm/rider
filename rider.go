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
	"time"
	"os"
	"path/filepath"
	"rider/logger"
	ctxt "context"
	"os/signal"
	"syscall"
	"strings"
	"rider/utils/file"
	"errors"
)

const (
	addr                    = ":8000"
	readTimeout             = 10 * time.Second
	writerTimeout           = 60 * time.Second
	maxHeaderBytes          = 1 << 20 //1MB
	defaultMultipartBodySze = 32 << 20
	ENV_Production          = "production"
	ENV_Development         = "development"
	ENV_Debug               = "debug"
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
	GetServer() *HttpServer //获取
	AddMiddleware(handlers ...HandlerFunc)
	Error(errorHandle func(c *Context, err string, code int))
}

//http服务的入口，用户初始化和缓存服务的一些信息
type rider struct {
	server  *HttpServer //注册服务用的serveMu，全局统一
	routers *Router
	appServer *http.Server
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

func SetEnvProduction() {
	SetEnvMode("production")
}

func SetEnvDevelopment() {
	SetEnvMode("development")
}

func SetEnvDebug() {
	SetEnvMode("debug")
}


//初始化服务入口组建
func New() *rider {
	server := newHttpServer()
	app := &rider{
		server:  server,
		routers: newRootRouter(server),
	}
	app.appServer = &http.Server{Handler: app.routers}
	//默认日志等级5 consoleLevel
	//日志会默认初始化，调用app.Logger(int)可以改变日志的输出等级
	app.server.logger = logger.NewLogger()
	app.server.logger.SetLevel(5)
	return app
}

func newRootRouter(server *HttpServer) *Router {
	_router := NewRouter()
	_router.isRoot = true
	_router.fullPath = "/"
	_router.Method = "ANY"
	_router.SetServer(server)
	return _router
}




//提供端口监听服务，监听rider里面的serveMux,调用http自带的服务启用方法
func (r *rider) Listen(port string) (err error) {
	if port == "" {
		port = addr
	}

	r.appServer.Addr = port
	r.appServer.ReadTimeout = readTimeout
	r.appServer.WriteTimeout = writerTimeout
	r.appServer.MaxHeaderBytes = maxHeaderBytes
	err = r.appServer.ListenAndServe()
	if err != nil {
		r.server.logger.ERROR(err.Error())
	}
	return
}


func (r *rider) Graceful(port string) {
	ch := make(chan os.Signal)
	var err error
	go func() {
		err = r.Listen(port)
		if err != nil {
			ch <- syscall.SIGINT
		}
	}()
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	r.appServer.Shutdown(ctxt.Background())
	time.Sleep(300 * time.Microsecond)
}

func (r *rider) Routers() *Router {
	return r.routers
}


//http请求的方法的入口（ANY, GET, POST...VIA）
//path：一个跟路径，函数内部根据这个根路径创建一个根路由routers，用来管理router子路由
//router：这个根路径对应的子路由入口。
func (r *rider) addRiderRouter(method string, path string, handlers ...HandlerFunc) {
	//将app的中间处理函数传给routers根路由(后面再由routers传给其各个子路由)
	//r.routers.FrontMiddleware(r.Middleware...)
	//将服务注入这个组建根路由，确保路由是创立在这个服务上
	//r.routers.SetServer(r.server)
	//将服务注入这个组建子路由，确保路由是创立在这个服务上
	//router.SetServer(r.server)
	//服务注册入口内部的入口
	switch method {
	case "ANY":
		r.routers.ANY(path, handlers...)
	case "GET":
		r.routers.GET(path, handlers...)
	case "POST":
		r.routers.POST(path, handlers...)
	case "HEAD":
		r.routers.HEAD(path, handlers...)
	case "DELETE":
		r.routers.DELETE(path, handlers...)
	case "PUT":
		r.routers.PUT(path, handlers...)
	case "PATCH":
		r.routers.PATCH(path, handlers...)
	case "OPTIONS":
		r.routers.OPTIONS(path, handlers...)
	case "CONNECT":
		r.routers.CONNECT(path, handlers...)
	case "TRACE":
		r.routers.TRACE(path, handlers...)
	}
}

//router：这个根路径对应的子路由入口。
func (r *rider) addChildRouter(path string, router ...IsRouterHandler) {

	if len(router) == 0 {
		r.server.logger.FATAL("nothing todo with no handler router")
	}

	for k, _ := range router {

		if k < len(router) - 1 {
			if _, ok := router[k].(HandlerFunc); !ok {
				r.server.logger.FATAL("put kid router at the end when register kid router")
			}
		} else {
			if route, ok := router[k].(*Router); !ok {
				r.server.logger.FATAL("if you want register handler directly do not use kid")
			} else {
				//将服务注入这个组建子路由，确保路由是创立在这个服务上
				route.SetServer(r.server)
			}
		}
	}

	//服务注册入口内部的入口
	r.routers.Kid(path, router...)
}

func (r *rider) Kid(path string, middleware ...IsRouterHandler) {
	r.addChildRouter(path, middleware...)
}

func (r *rider) ANY(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("ANY", path, handlers...)
}

func (r *rider) GET(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("GET", path, handlers...)
}

func (r *rider) POST(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("POST", path, handlers...)
}

func (r *rider) HEAD(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("HEAD", path, handlers...)
}

func (r *rider) OPTIONS(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("OPTIONS", path, handlers...)
}

func (r *rider) DELETE(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("DELETE", path, handlers...)
}

func (r *rider) PUT(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("PUT", path, handlers...)
}

func (r *rider) PATCH(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("PATCH", path, handlers...)
}

func (r *rider) TRACE(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("TRACE", path, handlers...)
}

func (r *rider) CONNECT(path string, handlers ...HandlerFunc) {
	r.addRiderRouter("CONNECT", path, handlers...)
}

//返回服务内部的http服务入口
func (r *rider) GetServer() *HttpServer {
	return r.server
}

//为app服务添加中间处理
func (r *rider) USE(handlers ...HandlerFunc) {
	r.routers.Middleware = append(r.routers.Middleware, handlers...)
}

//重写错误处理
func (r *rider) Error(errorHandle func(c Context, err string, code int)) {
	ErrorHandle = errorHandle
}

//设置模板路径（默认不缓存）
//tplDir以"/"开头，不会对其进行操作。如果直接以路径开头的，前面会默认跟上当前工作路径
func (r *rider) SetViews(tplDir string, extName string) (*render, error) {
	if !(strings.HasPrefix(tplDir, "/")) {
		tplDir = filepath.Join(file.GetCWD(), tplDir)
	}
	r.GetServer().tplDir = tplDir
	r.GetServer().tplExtName = extName
	if tplsRender, ok := r.GetServer().tplsRender.(*render); ok {
		appRender := tplsRender.registerTpl(tplDir, extName, "")
		return appRender, nil
	} else {
		return nil, errors.New("render is not implement BaseRender")
	}
}

//设置模板接口 (实现BaseRender接口的Render方法)
//须在SetViews方法之前调用
func (r *rider) ViewEngine(render BaseRender) {
	r.GetServer().tplsRender = render
}

//设置静态文件目录
func (r *rider) SetStatic(staticPath string) {
	if !(strings.HasPrefix(staticPath, "/")) {
		staticPath = filepath.Join(file.GetCWD(), staticPath)
	}
	f, err := os.Stat(staticPath)
	if err != nil {
		r.server.logger.FATAL(err.Error())
		return
	}
	if !f.IsDir() {
		r.server.logger.FATAL(staticPath + "不是路径，静态文件路径必须为目录")
		return
	}
	r.GET("/assets/(.*)", func(c Context) {
		c.SendFile(filepath.Join(staticPath, c.PathParams()[0]))
	})
}

//引入日志模块
func (r *rider) Logger(level int) *logger.LogQueue {
	//r.server.logger = logger.NewLogger()
	r.server.logger.SetLevel(level)
	return r.server.logger
}

//获取日志
func (r *rider) GetLogger() *logger.LogQueue {
	return r.server.logger
}
