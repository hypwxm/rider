package riderServer

import (
	"net/http"
	"sync"
	"fmt"
)

type HttpServer struct {
	pool      pool
	ServerMux *http.ServeMux
}

type pool struct {
	request  *sync.Pool
	response *sync.Pool
}

type Response struct {
	writer http.ResponseWriter
	Status int
	Size   int64
	body   []byte
	header http.Header
	sent   bool
}

type Request struct {
	*http.Request
	Query  map[string]string
	Body   map[string]string
	Method string
}

func (hs *HttpServer) Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}
