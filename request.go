package rider

import (
	"net/http"
	"net/url"
	"rider/cryptos"
	"io/ioutil"
	"strings"
	"rider/riderServer"
)


type Request struct {
	request *http.Request
	postBody   []byte
	isReadBody bool
	requestID  string
}

func NewRequest(r *http.Request) *Request {
	return (&Request{}).reset(r)
}

//reset response attr
func (req *Request) reset(r *http.Request) *Request {
	req.request = r
	req.isReadBody = false
	req.requestID = cryptos.GetUUID()
	return req
}

func (req *Request) release() {
	req.request = nil
	req.isReadBody = false
	req.postBody = nil
	req.requestID = ""
}

// RequestID get unique ID with current request
func (req *Request) RequestID() string {
	return req.requestID
}

// QueryStrings 返回查询字符串map表示
func (req *Request) QueryStrings() url.Values {
	return req.request.URL.Query()
}

/*
* 获取原始查询字符串
 */
func (req *Request) RawQuery() string {
	return req.request.URL.RawQuery
}

/*
* 根据指定key获取对应value
 */
func (req *Request) QueryString(key string) string {
	return req.request.URL.Query().Get(key)
}

func (req *Request) FormFile(key string) (*riderServer.UploadFile, error) {
	file, header, err := req.request.FormFile(key)
	if err != nil {
		return nil, err
	} else {
		return riderServer.NewUploadFile(file, header), nil
	}
}

/*
* 获取包括post、put和get内的值
 */
func (req *Request) FormValues() map[string][]string {
	req.parseForm()
	return map[string][]string(req.request.Form)
}

func (req *Request) parseForm() error {
	if strings.HasPrefix(req.QueryHeader("Content-Type"), "multipart/form-data") {
		if err := req.request.ParseMultipartForm(riderServer.DefaultMaxMemory); err != nil {
			return err
		}
	} else {
		if err := req.request.ParseForm(); err != nil {
			return err
		}
	}
	return nil
}

func (req *Request) ContentType() string {
	return req.request.Header.Get("Content-Type")
}

func (req *Request) QueryHeader(key string) string {
	return req.request.Header.Get(key)
}

//Deprecated: Use the PostFormValue instead
//returns the first value for the named component of the POST
// or PUT request body. URL query parameters are ignored.
func (req *Request) PostString(key string) string {
	return req.request.PostFormValue(key)
}

/*
* 获取post提交的字节数组
 */
func (req *Request) PostBody() []byte {
	if !req.isReadBody {
		bts, err := ioutil.ReadAll(req.request.Body)
		if err != nil {
			return []byte{}
		} else {
			req.isReadBody = true
			req.postBody = bts
		}
	}
	return req.postBody
}

//RemoteAddr to an "IP" address
func (req *Request) RemoteIP() string {
	fullIp := req.request.RemoteAddr
	//special: if run in win10, localIp will be like "[::]:port"
	//fixed for #20 cann't get RemoteIP and RemoteAddr in win10
	lastFlagIndex := strings.LastIndex(fullIp, ":")
	if lastFlagIndex >= 0 {
		return fullIp[:lastFlagIndex]
	}
	return fullIp
}

//RemoteAddr to an "IP:port" address
func (req *Request) FullRemoteIP() string {
	fullIp := req.request.RemoteAddr
	return fullIp
}

// Path returns requested path.
//
// The path is valid until returning from RequestHandler.
func (req *Request) Path() string {
	return req.request.URL.Path
}

// IsAJAX returns if it is a ajax request
func (req *Request) IsAJAX() bool {
	return req.request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (req *Request) Url() string {
	return req.request.URL.String()
}

//获取请求的httpMethod
func (req *Request) Method() string {
	return req.request.Method
}