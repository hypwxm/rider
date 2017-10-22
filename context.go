package rider

import (
	"net/http"
	"container/list"
	"errors"
	ctxt "context"
	"net/url"
	"os"
	"rider/logger"
	"strings"
	"encoding/json"
)

type Contexter interface {
	//NewContext(w http.ResponseWriter, r *http.Request) *Context                            //初始化一个Context
	Next(err ...error) error                                    //实现中间件的链表处理
	setCurrent(element *list.Element)                           //设置链表中当前在处理处理器
	setStartHandler() *list.Element                             //设置开始处理的第一个中间件
	getCurrentHandler() (HandlerFunc, error)                    //获取当前处理中的中间件
	startHandleList() error                                     //开始处理中间件
	release()                                                   //释放Context
	reset(w *Response, r *Request, server *HttpServer) *Context //初始化Context

	//处理http.request部分
	SetLocals(key string, value interface{})                           //给context传递变量，该变量在整个请求的传递中一直有效
	GetLocals(key string) interface{}                                  //通过SetLocals设置的值可以在下个中间件中获取
	Query() url.Values                                                 //只回去请求url中的查询字符串querystring的map[]string值
	QueryString(key string) string                                     //根据字段名直接查询querystring某个字段名对应的值
	Body() url.Values                                                  //只获取请求体内的请求参数，
	BodyValue(key string) string                                       //根据字段名直接查询"请求体"某中个字段名对应的值
	Params() map[string]string                                         //获取请求路由 /:id/:xx中值
	Param(key string) string                                           //获取请求路由中某字段的值
	FormFile(key string) (*UploadFile, error)                          //当请求头的content-type为multipart/form-data时，获取请求中key对应的文件信息(多个文件时，只会获取第一个)
	FormFiles(key string) ([]*UploadFile, error)                       //当请求头的content-type为multipart/form-data时，获取请求中key对应的文件列表
	StoreFormFile(file *UploadFile, fileName string) (int64, error)    //将file保存，fileName指定完整的路径和名称（先调用FormFile或者FormFiles将返回的file传入即可）
	StoreFormFiles(files []*UploadFile, path string) ([]string, error) //先通过FormFiles获取文件列表，指定path目录，存储文件的文件夹。文件名将会用随机字符串加文件的后缀名（file.GetFileExt()）
	Header() http.Header                                               //获取请求头信息
	HeaderValue(key string) string                                     //根据key获取请求头某一字段的值
	URL() string                                                       //获取请求头的完整url
	Path() string                                                      //获取请求头的path
	Method() string                                                    //获取响应头的HTTP方法
	ClientIP() string                                                  //获取请求来源的IP地址
	RequestID() string                                                 //获取分配的请求id
	IsAjax() bool                                                      //判断请求是否为ajax
	Cookies() []*http.Cookie                                           //获取请求的cookies
	CookieValue(key string) (string, error)                            //获取请求体中某一字段的cookie值
	Status() int  //获取响应的状态码
	Hijack() (*HijackUp, error) //将http请求升级为hijack，hijack的信息保存在HijackUp中
	SendFile(path string) error
	PathParams() []string //通过正则匹配后得到的路径上的一些参数
	//Response
	//Responser
	ResponseHeader() http.Header           //获取完整的响应头信息
	ResponseHeaderValue(key string) string //获取响应头的某一字段值
	RemoveCookie(cookieName string)        //删除cookie

	//模板
	Render(tplName string, data interface{})                 //负责模板渲染
	Download(fileName string, name string, typ string) error //下载，fileName为完整路径,name为下载时指定的下载名称，传""使用文件本身名称
}

type Responser interface {
	SetHeader(key, value string)             //设置响应头
	AddHeader(key, value string)             //给响应头添加值
	SetCType(contentType string)             //给响应头设置contenttype
	WriteHeader(code int)                    //设置响应头的状态码
	Header() http.Header                     //返回响应头信息
	HeaderValue(key string) string           //返回某一字段的响应头信息
	Redirect(code int, targetUrl string)     //重定向
	Write(data []byte) (size int, err error) //给客户端发送响应,返回发送的字节长度
	SetCookie(cookie http.Cookie)            //设置cookie
	//Hijack() (*HijackUp, error)  //将response升级为hijack，升级后，Send,SendJson等的相应方法也会调用hijack相关的
}

var (
	_ Responser = &Response{}
	_ Contexter = &Context{}
)

// "##"表示从pool取得context时必须初始化的值
type Context struct {
	request             *Request        //##
	response            *Response       //##
	handlerList         *list.List      //##
	currentHandler      *list.Element   //##
	ctx                 ctxt.Context    //##标准包的context
	cancel              ctxt.CancelFunc //##
	isHijack            bool            //##
	hijacker            *HijackUp       //##
	isEnd               bool            //##
	done                chan int        //##
	server              *HttpServer     //整个服务引用的server都是同一个
	ended               chan int        //##
	committed           bool            //##表示状态码是否已经发送（writeHeader有无调用）
	Logger              *logger.LogQueue
	Jwt                 *riderJwter //用于存储jwt
}

func newContext(w http.ResponseWriter, r *http.Request, server *HttpServer) *Context {
	//从pool取得context
	context := basePool.context.Get().(*Context)
	//从pool取得requset（注意初始化）
	request := NewRequest(r)
	//从pool取得response（注意初始化）
	response := NewResponse(w, server)
	//（初始化context）
	context.reset(response, request, server)
	return context
}

func releaseContext(c *Context) {
	//响应结束，释放request，response，context回poo，方便其他请求取得（释放时记得参数初始化）
	c.response.release()
	basePool.response.Put(c.response)
	c.request.release()
	basePool.request.Put(c.request)
	c.release()
	basePool.context.Put(c)
}

func (c *Context) reset(w *Response, r *Request, server *HttpServer) *Context {
	c.request = r
	c.response = w
	c.currentHandler = nil
	c.handlerList = list.New()
	//c.ctx = ctxt.Background()
	c.ctx = c.request.request.Context()
	c.isHijack = false
	c.hijacker = nil
	//c.isEnd = false
	c.done = make(chan int)
	c.ended = make(chan int, 1)
	c.server = server
	c.committed = false
	c.Logger = server.logger
	c.Jwt = nil
	//c.ctx, c.cancel = ctxt.WithTimeout(ctxt.Background(), writerTimeout)
	//go c.timeout()
	return c
}

//释放context里面的信息
func (c *Context) release() {
	c.response = nil
	c.request = nil
	c.currentHandler = nil
	c.handlerList = list.New()
	c.ctx = nil
	c.isHijack = false
	c.hijacker = nil
	c.isEnd = false
	c.done = nil
	//c.ended = nil
	c.hijacker = nil
	c.Logger = nil
	c.Jwt = nil
	c.committed = false
}

//处理下一个中间件
func (c *Context) Next(err ...error) error {
	if c.currentHandler == nil {
		//未知错误
		return errors.New("unknown error of nil handler")
	}

	if len(err) > 0 && err[0] != nil {
		HttpError(c, err[0].Error(), http.StatusInternalServerError)
		return nil
	}

	next := c.currentHandler.Next()
	if next == nil {
		//处理器已处理完毕，或者一些位置错误
		return errors.New("all handler was done")
	}

	//先设置current
	c.setCurrent(next)
	//在执行链表处理函数
	next.Value.(HandlerFunc)(c)
	return nil
}

//设置链表中当前在处理处理器
func (c *Context) setCurrent(e *list.Element) {
	c.currentHandler = e
}

//设置链表中开始处理是的起始handler
func (c *Context) setStartHandler() *list.Element {
	e := c.handlerList.Front()
	c.currentHandler = e
	return e
}

//获取当前正在处理中的handler
func (c *Context) getCurrentHandler() (HandlerFunc, error) {
	if c.currentHandler != nil {
		return c.currentHandler.Value.(HandlerFunc), nil
	}
	return nil, errors.New("处理函数为nil")
}

//开始处理链表的中的处理器
func (c *Context) startHandleList() error {
	c.setCurrent(c.setStartHandler())
	//开始处理
	handlerFunc, err := c.getCurrentHandler()
	if err != nil {
		return err
	}
	handlerFunc(c)
	return nil
}

//往context中传入临时信息
func (c *Context) SetLocals(key string, val interface{}) {
	c.ctx = ctxt.WithValue(c.ctx, key, val)
}

//从context中获取临时信息
func (c *Context) GetLocals(key string) interface{} {
	return c.ctx.Value(key)
}

//获取request的query  map[string][]string
func (c *Context) Query() url.Values {
	return c.request.Query()
}

//根据key查询query
func (c *Context) QueryString(key string) string {
	return c.request.QueryValue(key)
}

//获取request请求体  map[string][]string
func (c *Context) Body() url.Values {
	return c.request.Body()
}

//根据key查询body里面的某一字段的第一个值
func (c *Context) BodyValue(key string) string {
	return c.request.BodyValue(key)
}

//根据客户端传过来的字段名会返回第一个文件
func (c *Context) FormFile(key string) (*UploadFile, error) {
	return c.request.FormFile(key)
}

//根据客户端传过来的字段名返回文件列表
func (c *Context) FormFiles(key string) ([]*UploadFile, error) {
	return c.request.FormFiles(key)
}

//保存客服端传过来的文件（multipart/form-data）
func (c *Context) StoreFormFile(file *UploadFile, fileName string) (int64, error) {
	return c.request.StoreFormFile(file, fileName)
}

//保存客户端传过来的文件列表（multipart/form-data）
func (c *Context) StoreFormFiles(files []*UploadFile, path string) ([]string, error) {
	return c.request.StoreFormFiles(files, path)
}

//获取path匹配的所有参数/:id/:xx
func (c *Context) Params() map[string]string {
	return c.request.Params()
}

//获取path匹配的参数/:id/:xx
func (c *Context) Param(key string) string {
	return c.request.Param(key)
}

//获取正则匹配路径后的一些参数
func (c *Context) PathParams() []string {
	return c.request.pathParams
}

//获取requestID
func (c *Context) RequestID() string {
	return c.request.requestID
}

//获取请求头
func (c *Context) Header() http.Header {
	return c.request.Header()
}

//获取请求头的某一字段的值
func (c *Context) HeaderValue(key string) string {
	return c.request.HeaderValue(key)
}

//获取请求的URL
func (c *Context) URL() string {
	return c.request.Url()
}

//获取请求的PATH
func (c *Context) Path() string {
	return c.request.Path()
}

//获取请求的HTTP方法
func (c *Context) Method() string {
	return c.request.Method()
}

//RemoteAddr to an "IP" address
func (c *Context) ClientIP() string {
	return c.request.ClientIP()
}

//判断请求是否为ajax
func (c *Context) IsAjax() bool {
	return c.request.IsAJAX()
}

//获取请求的cookies
func (c *Context) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

//获取某一字段的cookie
func (c *Context) CookieValue(key string) (string, error) {
	return c.request.CookieValue(key)
}

//设置响应头
func (c *Context) SetHeader(key, value string) {
	if c.isHijack {
		c.hijacker.SetHeader(key, value)
	} else {
		c.response.SetHeader(key, value)
	}
}

//添加响应头
func (c *Context) AddHeader(key, value string) {
	if c.isHijack {
		c.hijacker.AddHeader(key, value)
	} else {
		c.response.AddHeader(key, value)
	}
}

//设置响应的contenttype
func (c *Context) SetCType(contentType string) {
	if c.isHijack {
		c.hijacker.SetHeader("Content-Type", contentType)
	} else {
		c.response.SetCType(contentType)
	}
}


//获取响应头信息
func (c *Context) ResponseHeader() http.Header {
	if c.isHijack {
		return c.hijacker.Header()
	} else {
		return c.response.Header()
	}
}

//获取响应头的某一字段的值
func (c *Context) ResponseHeaderValue(key string) string {
	if c.isHijack {
		return c.hijacker.HeaderValue(key)
	} else {
		return c.response.HeaderValue(key)
	}
}

//设置cookies
func (c *Context) SetCookie(cookie http.Cookie) {
	if cookie.Path == "" {
		cookie.Path = "/"
	}
	if c.isHijack {
		c.hijacker.SetCookie(cookie)
	} else {
		c.response.SetCookie(cookie)
	}
}

//删除cookie
func (c *Context) RemoveCookie(cookieName string) {
	cookie := http.Cookie{Name: cookieName, MaxAge: -1, Path: "/"}
	c.SetCookie(cookie)
}

//
func (c *Context) CloseNotify() <-chan bool {
	return c.response.CloseNotify()
}



//获取响应状态码
func (c *Context) Status() int {
	if c.isHijack {
		return c.hijacker.status
	} else {
		return c.response.status
	}
}


func (c *Context) writeHeader(code int) {
	if c.isHijack {
		c.hijacker.WriteHeader(code)
	} else {
		c.response.WriteHeader(code)
	}
}


//给客户端发送响应
func (c *Context) Send(code int, data []byte) (size int, err error) {
	if c.ResponseHeader().Get(HeaderContentType) == "" {
		c.SetHeader(HeaderContentType, http.DetectContentType(data))
	}

	c.writeHeader(code)
	if c.isHijack {
		return c.hijacker.Write(data)
	} else {
		return c.response.Write(data)
	}
}

//发送json格式数据给客户端
func (c *Context) SendJson(code int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	c.SetCType("application/json")
	c.Send(code, jsonData)
}

//重定向
func (c *Context) Redirect(code int, targetUrl string) {
	if code < 300 || code > 308 {
		c.server.logger.PANIC("Invalid redirect status code")
		return
	}
	if code == 301 {
		if c.Method() != http.MethodGet {
			code = 307
		}
	}
	c.response.Redirect(code, targetUrl)
}

//hijack相关实现

//升级responsewrite为hijack
func (c *Context) Hijack() (*HijackUp, error) {
	var originHeader http.Header = c.ResponseHeader()
	hj, ok := c.response.writer.(http.Hijacker)
	if !ok {
		return nil, errors.New("serve can not upgrade to hijack")
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		return nil, errors.New("Hijack error: " + err.Error())
	}
	c.isHijack = true
	c.hijacker = &HijackUp{conn: conn, bufrw: bufrw, header: originHeader}
	return c.hijacker, nil
}

//模板渲染
func (c *Context) Render(tplName string, data interface{}) {
	var err error
	if c.isHijack {
		c.hijacker.setHeaders()
		err = c.server.tplsRender.Render(c.hijacker.bufrw, tplName, data)
		err = c.hijacker.bufrw.Flush()
		if err != nil {
			HttpError(c, err.Error(), 500)
			return
		}
		c.hijacker.Close()
	} else {
		err = c.server.tplsRender.Render(c.response.writer, tplName, data)
	}
	if err != nil {
		HttpError(c, err.Error(), 404)
		return
	}
}

//文件服务器
func (c *Context) SendFile(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		HttpError(c, err.Error(), 404)
		return err
	}
	defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
		HttpError(c, err.Error(), 404)
		return err
	}
	if fi.IsDir() {
		HttpError(c, "file is invalid type directory", 404)
		return errors.New("file is invalid type directory")
	}

	c.response.Size = fi.Size()
	if setWeakEtag(c, fp, c.request.request) {
		c.Send(304, []byte(""))
		return nil
	}

	if c.isHijack {
		http.ServeContent(c.hijacker, c.request.request, fi.Name(), fi.ModTime(), fp)
	} else {
		http.ServeContent(c.response.writer, c.request.request, fi.Name(), fi.ModTime(), fp)
	}
	return nil
}

//文件下载
//filename为文件完整路径
//name指定文件的下载名称，为空时使用filename解析出文件名
//typ是下载的文件的模式inline和attachment（inline浏览器会以能够打开文件的软件默认打开下载文件，如果文件格式不支持打开，则表现和attachment一样只是下载）
func (c *Context) Download(fileName string, name string, typ string) error {
	if strings.TrimSpace(fileName) == "" {
		//filename不能为空
		c.server.logger.WARNING(fileName, "filename can not be empty. ")
		return errors.New("filename can not be empty. ")
	}
	if strings.TrimSpace(typ) == "" {
		typ = "attachment"
	}
	c.SetHeader("Content-Disposition", typ+";filename="+name)
	c.SendFile(fileName)
	return nil
}
