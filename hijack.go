package rider

import (
	"net"
	"bufio"
	"strings"
	"strconv"
	"net/http"
)

type HijackUp struct {
	conn net.Conn
	bufrw *bufio.ReadWriter
	header http.Header   //header引用了转换成hijack之前的response的header，
	status int
	Size int64
}

//添加默认的响应头
var defaultHijackHeader string = "HTTP/1.1 200 OK\r\n"

//设置响应头信息
//想要hijack正确的响应http请求，必须调用这个方法
func (hj *HijackUp) setHeaders() {
	if hj.header.Get("Content-Encoding") != "" {
		return
	}
	var header string = defaultHijackHeader
	for k, v := range hj.header {
		//如果响应头是setcookie，要一行行分开写
		if k == "Set-Cookie" {
			for _, cv := range v {
				header += k + ": " + cv + "\r\n"
			}
			continue
		}

		header += k + ": " + strings.Join(v, ";") + "\r\n"
	}
	hj.bufrw.WriteString(header + "\r\n")
}

//给hijack添加响应头
func (hj *HijackUp) AddHeader(key string, value string) {
	hj.Header()
	hj.header[key] = append(hj.header[key], value)
}

func (hj *HijackUp) SetHeader(key string, value string) {
	hj.Header()
	hj.header[key] = []string{value}
}

//获取响应头信息
func (hj *HijackUp) Header() http.Header {
	if hj.header == nil {
		hj.header = make(http.Header)
	}
	return hj.header
}

//获取响应头某一字段的值
func (hj *HijackUp) HeaderValue(key string) string {
	if val, ok := hj.header[key]; ok {
		if len(val) > 0 {
			return val[0]
		}
	}
	return ""
}

//获取响应头的http格式
func (hj *HijackUp) GetHttp() string {
	return "HTTP/1.1"
}

//关闭hijack conn
func (hj *HijackUp) Close() {
	hj.conn.Close()
}

//设置响应头状态码
//hijack的状态码重复设置以最后一次设置为发送的状态码
// 实现http.ResponseWriter
// 当code为304的时候会添加一个Last-Modified响应头到header，不会再发送实体部分，所以不会走到send方法
// 由于hijack的响应头需要自己手动发送，所所以当code为304的时候要主动设置响应头，并且flush到响应中去，还得close conn，否则，连接会一直挂起。

func (hj *HijackUp) WriteHeader(code int) {
	hj.status = code
	defaultHijackHeader = "HTTP/1.1 " + strconv.Itoa(code) + " " + http.StatusText(code) + "\r\n"
	if code == http.StatusNotModified {
		//hj.Send([]byte("xx"))
		hj.setHeaders()
		hj.bufrw.Flush()
		hj.Close()
	}
}

//给客户端发送数据
func (hj *HijackUp) Write(data []byte) (size int, err error) {
	hj.setHeaders()
	size = 0
	size, err = hj.bufrw.Write(data)
	hj.Size += int64(size)
	if err == nil {
		err = hj.bufrw.Flush()
	}
	hj.Close()
	return
}

//设置cookie
func (hj *HijackUp) SetCookie(cookie http.Cookie) {
	if v := (&cookie).String(); v != "" {
		hj.AddHeader("Set-Cookie", v)
	}
	//无需实行，，因为写cookie的方法时http.SetCookie(w, &http.Cookie{}) 。w实现了io.writer接口。hijack也实现了这个接口，所以直接用http.SetCookie就可以了
}

//设置contenttype
func (hj *HijackUp) SetCType(contentType string) {
	hj.SetHeader("Content-Type", contentType)
}
