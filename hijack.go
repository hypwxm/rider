package rider

import (
	"net"
	"bufio"
	"strings"
	"encoding/json"
	"strconv"
	"net/http"
)

type HijackUp struct {
	conn net.Conn
	bufrw *bufio.ReadWriter
	header map[string][]string
}

var defaultHijackHeader string = "HTTP/1.1 200 OK\r\n"

//给客户端发送数据
func (hj *HijackUp) Send(data []byte) (size int) {
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
	size, err := hj.bufrw.Write(data)
	if err == nil {
		err = hj.bufrw.Flush()
		hj.Close()
		//time.Sleep(10E9)
	} else {
		panic(err)
	}
	return
}

//发送json格式的数据给客户端
func (hj *HijackUp) SendJson(data interface{}) (size int) {
	hj.SetHeader("Content-Type", "application/json;charset=utf-8")
	dataJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	hj.Send(dataJson)
	hj.Close()
	return len(dataJson)
}

//给hijack添加响应头
func (hj *HijackUp) AddHeader(key string, value string) {
	hj.header[key] = append(hj.header[key], value)
}

func (hj *HijackUp) SetHeader(key string, value string) {
	hj.header[key] = []string{value}
}

//获取响应头信息
func (hj *HijackUp) GetHeader() map[string][]string {
	return hj.header
}

//获取响应头某一字段的值
func (hj *HijackUp) GetHeaderValue(key string) string {
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
func (hj *HijackUp) SetStatusCode(code int) {
	defaultHijackHeader = "HTTP/1.1 " + strconv.Itoa(code) + " " + http.StatusText(code) + "\r\n"
}

//设置cookie
func (hj *HijackUp) SetCookie(cookie http.Cookie) {

}