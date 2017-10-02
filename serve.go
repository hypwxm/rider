package rider

import (
	"net/http"
	"sync"
	"runtime/debug"
	"log"
	"errors"
)

type HttpServer struct {
	ServerMux *http.ServeMux
}

type pool struct {
	request  *sync.Pool
	response *sync.Pool
	context  *sync.Pool
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

func HttpError(c *Context, err string, code int) error {
	if c.isEnd {
		return errors.New("response sent again")
	}
	c.End()
	hijacker, errh := c.Hijack()
	if errh != nil {
		return errh
	}
	hijacker.SetStatusCode(code)
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
		log.Println(err)
		log.Println(string(debug.Stack()))
		errMsg.Stack = string(debug.Stack())
	}
	hijacker.SendJson(errMsg)
	return nil
}
