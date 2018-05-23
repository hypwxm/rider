package rider

import (
	"fmt"
	"net/http"
)

//实现一个http.ResponseWriter
type riderWriter struct{}

type wbuff struct {
	body []byte
}

var wf *wbuff = &wbuff{}

func (w *wbuff) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

func (rw *riderWriter) Header() http.Header {
	return http.Header{}
}

func (rw *riderWriter) Write(b []byte) (int, error) {
	fmt.Fprint(wf, string(b))
	return len(b), nil
}

func (rw *riderWriter) WriteHeader(i int) {
	fmt.Fprint(wf, i)
}
