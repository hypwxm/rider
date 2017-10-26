package rider

import (
	"compress/gzip"
	"strings"
	"net/http"
	"sync"
	"io/ioutil"
	"path/filepath"
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

func Gzip(level int) HandlerFunc {
	var gzPool sync.Pool
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}
	return func(c Context)  {
		if !shouldCompress(c.Request().request) {
			c.Next()
			return
		}
		res := c.Response()
		originW := res.writer

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)
		gz.Reset(res.writer)
		c.SetHeader(HeaderContentEncoding, "gzip")
		c.AddHeader(HeaderVary, HeaderAcceptEncoding)
		gw := &gzipWriter{}
		gw.ResponseWriter = res.writer
		gw.writer = gz
		res.writer = gw

		defer func() {
			if res.Size == 0 {
				//当发生panic的时候会走到这一步；将response的writer恢复为原来的writer，因为gzip的writer无法正确处理错误
				if res.Header().Get(HeaderContentEncoding) == "gzip" {
					res.Header().Del(HeaderContentEncoding)
				}
				res.writer = originW
			}
			gz.Close()
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
	if !strings.Contains(req.Header.Get(HeaderAcceptEncoding), "gzip") {
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

/*
func (g *gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw := &bufio.ReadWriter{}
	writer := bufio.NewWriter(g)
	rw.Writer = writer
	n, _ := net.Dial("tcp", "mv.51mzzk.com:5000")
	return n, rw, nil
}
*/
