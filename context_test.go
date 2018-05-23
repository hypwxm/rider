package rider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/hypwxm/rider/logger"
)

func newTestContext() (Context, http.ResponseWriter, *http.Request) {
	var context Context
	var hw *riderWriter
	var req *http.Request
	hw = &riderWriter{}
	url, _ := url.Parse("localhost:5000/a/x")
	req = &http.Request{
		Method: http.MethodGet,
		URL:    url,
	}
	server := &HttpServer{}
	server.logger = logger.NewLogger()
	context = newContext(hw, req, server)
	return context, hw, req
}

func TestNewContext(t *testing.T) {
	context, hw, req := newTestContext()
	if context.Response().writer != hw {
		t.Error("设置response错误")
	}
	if context.Request().request != req {
		t.Error("设置request错误")
	}

}

func TestCSend(t *testing.T) {
	context, _, _ := newTestContext()

	context.Send(200, []byte("test"))
	//context.Send([]byte("test"))

	body := fmt.Sprintf("%s", wf.body)
	if body != "200test" {
		t.Error("context.response.send error")
	}
	wf.body = []byte{}

}

func TestCSendJson(t *testing.T) {
	context, _, _ := newTestContext()

	context.SendJson(200, map[string]string{
		"a": "a",
	})
	mapdata := make(map[string]string)

	t.Logf("%s", wf.body[3:]) //200{"a":"a"}
	err := json.Unmarshal(wf.body[3:], &mapdata)
	if err != nil {
		t.Error(err)
	}
	if mapdata["a"] != "a" {
		t.Error("context.response.send error")
	}

	header := context.HeaderValue("Content-Type")
	t.Logf("%s", header)
	wf.body = []byte{}
}

func TestHijackSend(t *testing.T) {
	context, _, _ := newTestContext()
	context.Hijack()

	context.Send(200, []byte("test"))
	t.Logf("%s", wf.body)
	body := fmt.Sprintf("%s", wf.body)
	if body != "200test" {
		t.Error("context.response.send error")
	}
	wf.body = []byte{}
}

func TestHijackSendJson(t *testing.T) {
	context, _, _ := newTestContext()

	context.Hijack()

	context.SendJson(200, map[string]string{
		"a": "a",
	})
	mapdata := make(map[string]string)

	t.Logf("%s", wf.body[3:]) //200{"a":"a"}
	err := json.Unmarshal(wf.body[3:], &mapdata)
	if err != nil {
		t.Error(err)
	}
	if mapdata["a"] != "a" {
		t.Error("context.response.send error")
	}

	header := context.HeaderValue("Content-Type")
	t.Logf("%s", header)
	wf.body = []byte{}
}

func BenchmarkSend(b *testing.B) {

}
