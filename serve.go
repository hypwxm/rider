package rider

import (
	"net/http"
	"sync"
	"runtime/debug"
	"rider/logger"
)

type HttpServer struct {
	ServerMux *http.ServeMux
	tplDir string
	tplExtName string
	tplsRender BaseRender
	logger *logger.LogQueue
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
		ServerMux: http.NewServeMux(),
		tplsRender: newRender(),
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
			return &Context{}
		},
	},
}

func (h *HttpServer) NewHttpServer() *HttpServer {
	return &HttpServer{}
}
var ErrorHandle func(c *Context, err string, code int)


func HttpError(c *Context, err string, code int) {
	if c.isEnd {
		c.server.logger.PANIC("can not send a response again")
		return
	}
	//c.End()

	//错误处理函数会将responsewriter转化为hijack；为了能够立即将响应返回给客户端
	//如果不转成hijack，将会导致错误处理会一直等待主程序其他的处理完成，但是其他处理完成的时候已经release所有变量，导致send时找不到对应的writer
	var hijacker *HijackUp
	var errh error
	if c.isHijack {
		hijacker = c.hijacker
	} else {
		hijacker, errh = c.Hijack()
	}

	if errh != nil {
		c.server.logger.PANIC(errh, "\r\n", string(debug.Stack()))
		return
	}
	hijacker.SetStatusCode(code)
	ErrorHandle(c, err, code)
}


//全局的错误处理，创建服务可以直接重写该方法
func init() {
	ErrorHandle = func(c *Context, err string, code int) {
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
			c.server.logger.DEBUG(err, "\r\n", string(debug.Stack()))
			errMsg.Stack = string(debug.Stack())
		}
		c.SendJson(errMsg)
	}
}
