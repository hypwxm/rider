package riderServer

import "net/http"

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		writer: w,
	}
}


func (r *Response) reset(w http.ResponseWriter) {
	r.writer = w
	r.header = w.Header()
	r.Status = http.StatusOK
	r.Size = 0
	r.sent = false
}