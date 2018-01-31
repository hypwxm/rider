package rider

import (
	"html/template"
	"net/http"
	"rider/logger"
	"runtime/debug"
	"sync"
)

type HttpServer struct {
	openRender bool
	tplDir     string
	tplExtName string
	funcMap    template.FuncMap
	tplsRender BaseRender
	logger     *logger.LogQueue
}

type ErrorHandler interface {
	ErrorHandle(c *Context, err string, code int)
}

type pool struct {
	request  *sync.Pool
	response *sync.Pool
	context  *sync.Pool
}

func newHttpServer() *HttpServer {
	return &HttpServer{
		tplsRender: newRender(),
		openRender: false,
	}
}

//全局的pool
var basePool *pool = &pool{
	response: &sync.Pool{
		New: func() interface{} {
			return &Response{}
		},
	},
	request: &sync.Pool{
		New: func() interface{} {
			return &Request{}
		},
	},
	context: &sync.Pool{
		New: func() interface{} {
			return &context{}
		},
	},
}

func (h *HttpServer) NewHttpServer() *HttpServer {
	return &HttpServer{}
}

var ErrorHandle func(c Context, err string, code int)

func HttpError(c Context, err string, code int) {
	ErrorHandle(c, err, code)
}

//全局的错误处理，创建服务可以直接重写该方法
func init() {
	ErrorHandle = func(c Context, err string, code int) {
		errMsg := &Error{
			StatusCode: code,
			StatusText: http.StatusText(code),
			Error:      err,
		}
		if GlobalENV == ENV_Production {
			errMsg.Stack = ""
		}
		if GlobalENV == ENV_Development {
			errMsg.Stack = string(debug.Stack())
		}
		if GlobalENV == ENV_Debug {
			c.Logger().DEBUG(err, "\r\n", string(debug.Stack()))
			errMsg.Stack = string(debug.Stack())
		}
		c.SendJson(code, errMsg)
	}
}
