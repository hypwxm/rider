package rider

import (
	"net/http"
	"net/url"
	"rider/cryptos"
	"io/ioutil"
	"strings"
	"rider/riderServer"
	"path/filepath"
	"errors"
)

type Requester interface {
	NewRequest(r *http.Request) *Request                                           //用传入的request生成一个Request
	reset(r *http.Request) *Request                                                //初始化一个requset
	release()                                                                      //重制Request，销毁变量，放入pool
	RequestID() string                                                             //获取请求的id
	Query() url.Values                                                             //获取请求url里面的参数部分(单个值查询请用QueryString)
	RawQuery() string                                                              //获取请求的原始字符串  a=1&b=2
	QueryString(key string) string                                                 //获取某个字段的url参数值
	FormFile(key string) (*riderServer.UploadFile, error)                          //请求头为multipart/form-data时获取key字段域的第一个文件
	FormFiles(key string) ([]*riderServer.UploadFile, error)                       //请求头为multipart/form-data时获取key字段域的文件列表
	StoreFormFile(file *riderServer.UploadFile, fileName string) (int64, error)    //FormFile返回一个文件，然后传入这里的第一个参数，保存到fileName文件
	StoreFormFiles(files []*riderServer.UploadFile, path string) ([]string, error) //FormFile返回文件列表，传入第一个参数，保存的文件名是随机的字符串，后缀名取文件上传时的本身后缀名
	FormValues() url.Values                                                        //获取url和body里面的参数。同名参数，body的值会排在前名
	Body() url.Values                                                              //获取请求体的参数，不包括文件，不包括url参数
	BodyValue(key string) string                                                   //根据key获取请求体内具体某个字段的值
	ContentType() string                                                           //获取请求体的content-type
	Header(key string) string                                                      //获取请求头信息
	PostBody() []byte                                                              //以流的形式解析请求体，调用PostBody后，r.Body()将无法在获取内容（反之也一样）
	RemoteIP() string                                                              //获取远程请求者的ip地址
	FullRemoteIP() string                                                          //获取远程请求者的ip:port地址
	Path() string                                                                  //获取请求的url路径
	IsAJAX() bool                                                                  //验证请求是否为ajax请求
	Url() string                                                                   //获取完整的请求路径  /x/x/a?a=x&c=s
	Method() string                                                                //获取请求的方法
	Params() map[string]string                                                     //获取请求的路径参数
	Param(key string) string                                                       //获取某个请求路径参数的值（参数定义时的字段请勿重复，后面的直接覆盖之前的值）
}

type Request struct {
	request    *http.Request
	postBody   []byte
	isReadBody bool
	requestID  string
	params     map[string]string
}

func NewRequest(r *http.Request) *Request {
	return (&Request{}).reset(r)
}

//reset response attr
func (req *Request) reset(r *http.Request) *Request {
	req.request = r
	req.isReadBody = false
	req.params = make(map[string]string, 4)
	req.requestID = cryptos.GetUUID()
	return req
}

func (req *Request) release() {
	req.request = nil
	req.isReadBody = false
	req.postBody = nil
	req.params = nil
	req.requestID = ""
}

// RequestID get unique ID with current request
func (req *Request) RequestID() string {
	return req.requestID
}

// QueryStrings 返回查询字符串map表示
func (req *Request) Query() url.Values {
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

//当客户端的请求头的content-type为multipart/form-data时获取请求体内的文件
//要获取请求传过来的文件，必须先调用r.parseMultipartForm（调用r.FormFile时会自动调用r.parseMultipartForm所以无需调用，只能获取第一个文件）
//只会返回第一个文件
func (req *Request) FormFile(key string) (*riderServer.UploadFile, error) {
	file, header, err := req.request.FormFile(key)
	if err != nil {
		return nil, err
	} else {
		return riderServer.NewUploadFile(file, header), nil
	}
}

//根据key返回key的文件列表
func (req *Request) FormFiles(key string) ([]*riderServer.UploadFile, error) {
	if req.request.MultipartForm == nil {
		err := req.request.ParseMultipartForm(defaultMultipartBodySze)
		if err != nil {
			return nil, err
		}
	}
	files := req.request.MultipartForm.File[key]
	rFiles := []*riderServer.UploadFile{}
	for i, fileHeader := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			return nil, err
		}
		rFiles = append(rFiles, riderServer.NewUploadFile(file, fileHeader))
	}
	return rFiles, nil
}

//获取客户端传过来的文件，并且指定保存路径进行保存
func (req *Request) StoreFormFile(file *riderServer.UploadFile, fileName string) (int64, error) {
	return file.SaveFile(fileName)
}

//文件名会根据用生成的随机字符串，加上源文件后缀名，path为文件保存的路径
//返回文件名列表
func (req *Request) StoreFormFiles(files []*riderServer.UploadFile, path string) ([]string, error) {
	fileNames := []string{}
	if path == "" {
		return fileNames, errors.New("path is not allowed empty")
	}
	for _, file := range files {
		fileName := cryptos.GetRandString(10) + "." + file.GetFileExt()
		_, err := req.StoreFormFile(file, filepath.Join(path, fileName))
		if err != nil {
			return fileNames, err
		}
		fileNames = append(fileNames, fileName)
	}
	return fileNames, nil
}

/*
* 获取包括post、put和get内的值
 */
func (req *Request) FormValues() url.Values {
	req.parseForm()
	return map[string][]string(req.request.Form)
}

func (req *Request) parseForm() error {
	if strings.HasPrefix(req.HeaderValue("Content-Type"), "multipart/form-data") {
		if err := req.request.ParseMultipartForm(defaultMultipartBodySze); err != nil {
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

func (req *Request) Header() map[string][]string {
	return req.request.Header
}

func (req *Request) HeaderValue(key string) string {
	return req.request.Header.Get(key)
}

//获取请求体body内容 url.Value，，
func (req *Request) Body() url.Values {
	req.parseForm()
	return req.request.PostForm
}

//获取body各字段的第一个值
func (req *Request) BodyValue(key string) string {
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

//获取请求的路由参数
func (req *Request) Params() map[string]string {
	return req.params
}

//根据参数名称获取路由参数
func (req *Request) Param(key string) string {
	if param, ok := req.params[key]; ok {
		return param
	}
	return ""
}

//获取cookies
func (req *Request) Cookies() []*http.Cookie {
	return req.request.Cookies()
}

//获取某一字段的cookie
func (req *Request) CookieValue(key string) (string, error) {
	cookieObj, err := req.request.Cookie(key)
	if err != nil {
		return "", err
	}
	return url.QueryUnescape(cookieObj.Value)
}

//给请求体添加cookie
func (req *Request) AddReqCookie(cookie http.Cookie)  {
	req.request.AddCookie(&cookie)
}