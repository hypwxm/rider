package rider

import (
	"net/http"
	"net"
	"bufio"
	"io"
	"errors"
	"encoding/json"
	"runtime/debug"
	"log"
	"time"
)


type (
	Response struct {
		writer    http.ResponseWriter
		Status    int
		Size      int64
		body      []byte
		committed bool
		header    http.Header
		isEnd     bool
		isHijack  bool
	}

	gzipResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

func NewResponse(w http.ResponseWriter) (r *Response) {
	return (&Response{}).reset(w)
}

func (r *Response) Header() http.Header {
	return r.header
}


func (r *Response) HeaderValue(key string) string {
	return r.Header().Get(key)
}

func (r *Response) Redirect(code int, targetUrl string) {
	r.Header().Set("Catch-Control", "no-cache")
	r.Header().Set("Location", targetUrl)
	r.SetStatusCode(code)
}

func (r *Response) SetWriter(w http.ResponseWriter) *Response {
	r.writer = w
	return r
}

func (r *Response) AddHeader(key, value string) {
	r.Header().Add(key, value)
}

func (r *Response) SetHeader(key, val string) {
	r.Header().Set(key, val)
}

func (r *Response) SetContentType(contenttype string) {
	r.SetHeader("Content-Type", contenttype)
}

// SetStatusCode sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
// used to send error codes.
func (r *Response) SetStatusCode(code int) {
	if r.committed {
		log.Panic(errors.New("response already set status"), "\n", string(debug.Stack()))
	}
	r.Status = code
	r.writer.WriteHeader(code)
	r.committed = true
}

//给client发送消息，json，text，html，xml...(不发送模板,模板请用render)
func (r *Response) Send(data []byte) (size int) {
	if !r.committed {
		r.SetStatusCode(http.StatusOK)
	}

	if r.isEnd {
		log.Panic(errors.New("sent again after res was sent"), "\n", string(debug.Stack()))
	}

	r.End()
	r.writer.Write(data)
	return len(data)
}

//发送json格式的数据
func (r *Response) SendJson(data interface{}) (size int) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	r.SetContentType("application/json;charset=utf-8")
	r.Send(jsonData)
	return len(jsonData)
}

//stop current response
func (r *Response) End() {
	r.isEnd = true
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *Response) Flush() {
	r.writer.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See https://golang.org/pkg/net/http/#Hijacker
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.writer.(http.Hijacker).Hijack()
}

//reset response attr
func (r *Response) reset(w http.ResponseWriter) *Response {
	r.writer = w
	r.header = w.Header()
	r.Status = http.StatusOK
	r.Size = 0
	r.committed = false
	r.isEnd = false
	r.body = []byte{}
	setBaseResHeader(r.Header())
	return r
}

//reset response attr
func (r *Response) release() {
	r.writer = nil
	r.header = nil
	r.Status = http.StatusOK
	r.Size = 0
	r.body = []byte{}
	r.isEnd = false
	r.committed = false
}


//设置cookies
func (r *Response) SetCookie(cookie http.Cookie) {
	http.SetCookie(r.writer, &cookie)
}



//设置一些基本的响应头信息
func setBaseResHeader(header http.Header) {
	header.Set("Server", "rider")
	header.Set("X-DNS-Prefetch-Control", "off")
	header.Set("X-Download-Options", "noopen")
	header.Set("X-Frame-Options", "SAMEORIGIN")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("Connection", "keep-alive")
	header.Set("Date", time.Now().Format(time.RFC822))
}