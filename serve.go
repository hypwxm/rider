package rider

import (
	"net/http"
	"sync"
	"fmt"
)

type HttpServer struct {
	ServerMux *http.ServeMux
}

type pool struct {
	request  *sync.Pool
	response *sync.Pool
	context *sync.Pool
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


func (h *HttpServer) Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}
