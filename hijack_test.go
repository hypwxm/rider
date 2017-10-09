package rider

import (
	"testing"
	"net/http"
	"net"
	"bufio"
	"os"
)

func TestAddHeader(t *testing.T) {
	hj := &HijackUp{}
	hj.AddHeader("test", "test")
	header := hj.Header()
	if len(header) != 1 || len(header["test"]) != 1 && header["test"][0] != "test" {
		t.Error("调用AddHeader()添加响应头信息有错")
	}

	hj.AddHeader("test", "test2")

	if len(header) != 1 || len(header["test"]) != 2 && header["test"][0] != "test" && header["test"][1] != "test2" {
		t.Error("调用AddHeader()添加响应头信息有错")
	}
}

func TestSetHeader(t *testing.T) {
	hj := &HijackUp{}
	hj.SetHeader("test", "test")
	header := hj.Header()
	if len(header) != 1 || len(header["test"]) != 1 && header["test"][0] != "test" {
		t.Error("调用SetHeader()设置响应头信息有错")
	}
	hj.SetHeader("test", "test2")
	if len(header) != 1 || len(header["test"]) != 1 && header["test"][0] != "test2" {
		t.Error("调用SetHeader()设置响应头信息有错")
	}
}

func TestHeader(t *testing.T) {
	hj := &HijackUp{}
	hj.Header()["test"] = []string{"test"}
	if hj.header["test"][0] != "test" {
		t.Error("获取响应头信息有错")
	}
}

func TestHeaderValue(t *testing.T) {
	hj := &HijackUp{}
	hj.AddHeader("test", "test")
	if hj.HeaderValue("test") != "test" {
		t.Error("获取某一字段的响应头信息有错")
	}
}


func TestSetStatusCode(t *testing.T) {
	hj := &HijackUp{}
	hj.SetStatusCode(404)
	if defaultHijackHeader != "HTTP/1.1 404 Not Found\r\n" {
		t.Error(defaultHijackHeader)
	}
}

func TestSetCookie(t *testing.T) {
	hj := &HijackUp{}
	hj.SetCookie(http.Cookie{Name:"test", Value: "test"})
	header := hj.Header()
	cookies := header["Set-Cookie"]
	if cookies[0] != "test=test" {
		t.Error("设置cookie有错")
	}
}

func TestSetContentType(t *testing.T) {
	hj := &HijackUp{}
	hj.SetContentType("application/json")
	header := hj.Header()
	if header["Content-Type"][0] != "application/json" {
		t.Error("Content-Type =", header["Content-Type"])
	}
	hj.SetContentType("application/xx")
	if header["Content-Type"][0] != "application/xx" {
		t.Error("Content-Type =", header["Content-Type"])
	}
}

type myHijack struct{}
func (mh *myHijack) Header() http.Header {
	return nil
}
func (mh *myHijack) Write([]byte) (int, error) {
	return 0, nil
}
func (mh *myHijack) WriteHeader(int) {
}
func (mh *myHijack) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	rw := &bufio.ReadWriter{}
	writer := bufio.NewWriter(os.Stdout)
	rw.Writer = writer
	n, _ := net.Dial("tcp", "localhost:8000")
	return n, rw, nil
}

func TestSend(t *testing.T) {
	var mh interface{} = &myHijack{}
	h, ok := mh.(http.Hijacker)
	if !ok {
		t.Error("无法转换成hijack")
	}
	conn, buf, err := h.Hijack()
	if err != nil {
		t.Error(err)
	}
	hj := &HijackUp{conn: conn, bufrw: buf}
	hj.SetStatusCode(200)
	hj.AddHeader("test", "test")
	size := hj.Send([]byte("xx"))
	if size != 2 {
		t.Error(size)
	}
}