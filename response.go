package rider

import (
	"net/http"
	"net"
	"bufio"
	"net/textproto"
	"strings"
	"crypto/md5"
	"os"
	"io/ioutil"
	"encoding/hex"
)

type (
	Response struct {
		writer    http.ResponseWriter
		Status    int
		header    http.Header
		committed bool //响应吗发送状态
		isEnd     bool
		isHijack  bool
		server    *HttpServer
		Size      int64
	}
)

func NewResponse(w http.ResponseWriter, server *HttpServer) (r *Response) {
	response := basePool.response.Get().(*Response)
	return response.load(w, server)
}

func (r *Response) Header() http.Header {
	return r.header
}

func (r *Response) HeaderValue(key string) string {
	return r.Header().Get(key)
}

//redirect会忽略setstatuscode设置的状态码
func (r *Response) Redirect(code int, targetUrl string) {
	r.Header().Set("Location", targetUrl)
	r.WriteHeader(code)
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

func (r *Response) SetCType(contenttype string) {
	r.SetHeader("Content-Type", contenttype+"; charset=utf-8")
}


//获取状态码
func (r *Response) GetStatusCode() int {
	return r.Status
}

func (r *Response) WriteHeader(code int) {
	if r.committed {
		if r.isEnd {
			r.server.logger.PANIC("can not send response status after sending a response")
		} else {
			r.server.logger.PANIC("can not set the status code again")
		}
		return
	}
	r.Status = code
	r.writer.WriteHeader(code)
	r.committed = true
}

//给client发送消息，json，text，html，xml...(不发送模板,模板请用render)
func (r *Response) Write(data []byte) (size int, err error) {
	if !r.committed {
		r.WriteHeader(http.StatusOK)
	}

	if r.isEnd {
		r.server.logger.PANIC("sent again after res was sent")
		return
	}

	r.End()
	r.Size = int64(len(data))
	return r.writer.Write(data)
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
func (r *Response) load(w http.ResponseWriter, server *HttpServer) *Response {
	r.writer = w
	r.header = w.Header()
	r.Status = http.StatusOK
	r.committed = false
	r.isEnd = false
	r.server = server
	setConfigHeaders(r.Header())
	return r
}

//reset response attr
func (r *Response) release() {
	r.writer = nil
	r.header = nil
	r.Status = http.StatusOK
	r.isEnd = false
	r.committed = false
}

//在客户端关闭连接但还未发送响应体时，关闭连接
func (r *Response) CloseNotify() <-chan bool {
	/*notify := r.writer.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		r.server.logger.WARNING("HTTP connection just closed.")
	}()*/
	return r.writer.(http.CloseNotifier).CloseNotify()
}

//设置cookies
func (r *Response) SetCookie(cookie http.Cookie) {
	http.SetCookie(r.writer, &cookie)
}

//验证weaketag 无论更新否都会返回etag，如果没更新，返回的etag和传入的是一样的，bool为true响应304
func weakEtag(fi *os.File, r *http.Request) (string, bool) {
	if !DefaultConfig.EnableWeakEtag {
		return "", false
	}
	chunk, err := ioutil.ReadAll(fi)
	if err != nil {
		return "", false
	}
	md5Chunk := md5.New()
	md5Chunk.Write(chunk)
	ms := md5Chunk.Sum(nil)
	mhex := hex.EncodeToString(ms)
	etagStr := string(mhex[:16])
	newEtag := `W/"` + etagStr + `"`

	//只有就的etag和新的etag完全匹配才会返回true，即需要返回304notModified
	//任何不瞒住etag的，都会重新生成一个etag，并返回给客户端
	ifNoneMatch := r.Header.Get("If-None-Match")
	ifNoneMatch = textproto.TrimString(ifNoneMatch)
	start := 0
	if strings.HasPrefix(ifNoneMatch, "W/") {
		start = 2
	}
	if len(ifNoneMatch[start:]) < 2 || ifNoneMatch[start] != '"' {
		return newEtag, false
	}
	// ETag is either W/"text" or "text".
	// See RFC 7232 2.3.
	for i := start + 1; i < len(ifNoneMatch); i++ {
		c := ifNoneMatch[i]
		switch {
		// Character values allowed in ETags.
		case c == 0x21 || c >= 0x23 && c <= 0x7E || c >= 0x80:
		case c == '"':
			ifNoneMatch = string(ifNoneMatch[start+1:i])
		default:
			return newEtag, false
		}
	}
	if etagStr == ifNoneMatch {
		return newEtag, true
	}
	return newEtag, false
}

//设置一些基本的响应头信息
func setConfigHeaders(header http.Header) {
	header.Set("Server", Server)
	if DefaultConfigHeaders.Date != "" {
		header.Set(HeaderDate, DefaultConfigHeaders.Date)
	}
	if DefaultConfigHeaders.XContentTypeOptions != "" {
		header.Set(HeaderXContentTypeOptions, DefaultConfigHeaders.XContentTypeOptions)
	}
	if DefaultConfigHeaders.XXSSProtection != "" {
		header.Set(HeaderXXSSProtection, DefaultConfigHeaders.XXSSProtection)
	}
	if DefaultConfigHeaders.XFrameOptions != "" {
		header.Set(HeaderXFrameOptions, DefaultConfigHeaders.XFrameOptions)
	}
	if DefaultConfigHeaders.XDownloadOptions != "" {
		header.Set(HeaderXDownloadOptions, DefaultConfigHeaders.XDownloadOptions)
	}
}

func setWeakEtag(c *Context, fi *os.File, r *http.Request) bool {
	etag, ifEqual := weakEtag(fi, r)
	if etag != "" {
		c.SetHeader(HeaderEtag, etag)
		if ifEqual {
			c.SetStatusCode(http.StatusNotModified)
			return true
		}
	}
	return false
}