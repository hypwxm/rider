package rider

import (
	"compress/gzip"
	"strings"
	"net/http"
	"sync"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"net"
	"bufio"
)


const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

type gzipWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func Gzip(level int) func(c *Context) {
	var gzPool sync.Pool
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}
	return func(c *Context)  {
		acceptEncoding := c.Header().Get(HeaderAcceptEncoding)
		if strings.Index(acceptEncoding, "gzip") < 0 {
			c.Next()
			return
		}
		res := c.response

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)
		gz.Reset(res.writer)

		c.SetHeader(HeaderContentEncoding, "gzip")
		c.AddHeader(HeaderVary, HeaderAcceptEncoding)
		gw := &gzipWriter{}
		gw.ResponseWriter = res.writer
		gw.writer = gz
		res.writer = gw

		//c.response.writer = gw
		defer func() {
			gz.Close()
			c.SetHeader("Content-Length", fmt.Sprint(c.response.Size))
		}()
		c.Next()
	}
}



func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) WriteHeader(code int) {
	if code == http.StatusNoContent {
		g.Header().Del(HeaderContentEncoding)
		g.Header().Del(HeaderContentLength)
	}
	g.ResponseWriter.WriteHeader(code)
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

func shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}
	extension := filepath.Ext(req.URL.Path)
	if len(extension) < 4 { // fast path
		return true
	}

	switch extension {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	default:
		return true
	}
}

func (g *gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw := &bufio.ReadWriter{}
	writer := bufio.NewWriter(g)
	rw.Writer = writer
	n, _ := net.Dial("tcp", "mv.51mzzk.com:5000")
	return n, rw, nil
}
