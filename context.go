package rider

import (
	"net/http"
	"container/list"
	"errors"
	ctxt "context"
	"net/url"
	"rider/riderServer"
	"time"
	"fmt"
)

type Contexter interface {
	//NewContext(w http.ResponseWriter, r *http.Request) *Context                            //初始化一个Context
	Next(err ...error) error                                                       //实现中间件的链表处理
	setCurrent(element *list.Element)                                              //设置链表中当前在处理处理器
	setStartHandler() *list.Element                                                //设置开始处理的第一个中间件
	getCurrentHandler() (HandlerFunc, error)                                       //获取当前处理中的中间件
	startHandleList() error                                                        //开始处理中间件
	release()                                                                      //释放Context
	reset(w *Response, r *Request) *Context                                        //初始化Context
	SetLocals(key string, value interface{})                                       //给context传递变量，该变量在整个请求的传递中一直有效
	GetLocals(key string) interface{}                                              //通过SetLocals设置的值可以在下个中间件中获取
	CancelResponse()                                                               //在响应未结束前，可以随时终止响应的处理
	Hijack() (*HijackUp, error)                                                    //将http请求升级为hijack，hijack的信息保存在HijackUp中
	Query() url.Values                                                             //只回去请求url中的查询字符串querystring的map[]string值
	QueryString(key string) string                                                 //根据字段名直接查询querystring某个字段名对应的值
	Body() url.Values                                                              //只获取请求体内的请求参数，
	BodyValue(key string) string                                                   //根据字段名直接查询"请求体"某中个字段名对应的值
	Params() map[string]string                                                     //获取请求路由 /:id/:xx中值
	Param(key string) string                                                       //获取请求路由中某字段的值
	FormFile(key string) (*riderServer.UploadFile, error)                          //当请求头的content-type为multipart/form-data时，获取请求中key对应的文件信息(多个文件时，只会获取第一个)
	FormFiles(key string) ([]*riderServer.UploadFile, error)                       //当请求头的content-type为multipart/form-data时，获取请求中key对应的文件列表
	StoreFormFile(file *riderServer.UploadFile, fileName string) (int64, error)    //将file保存，fileName指定完整的路径和名称（先调用FormFile或者FormFiles将返回的file传入即可）
	StoreFormFiles(files []*riderServer.UploadFile, path string) ([]string, error) //先通过FormFiles获取文件列表，指定path目录，存储文件的文件夹。文件名将会用随机字符串加文件的后缀名（file.GetFileExt()）
	Header() map[string][]string                                                   //获取请求头信息
	HeaderValue(key string) string                                                 //根据key获取请求头某一字段的值
	URL() string                                                                   //获取请求头的完整url
	Path() string                                                                  //获取请求头的path
	Method() string                                                                //获取响应头的HTTP方法
	RemoteIP() string                                                              //获取请求来源的IP地址
	FullRemoteIP() string                                                          //获取请求来源的完整IP:PORT
	RequestID() string                                                             //获取分配的请求id
	IsAjax() bool                                                                  //判断请求是否为ajax
	Cookies() []*http.Cookie                                                       //获取请求的cookies
	CookieValue(key string) (string, error)                                        //获取请求体中某一字段的cookie值

	End() //调用该方法表示响应已经结束

	//Response
	//Responser
	SendHeader() http.Header           //获取完整的响应头信息
	SendHeaderValue(key string) string //获取响应头的某一字段值
	SendCookie(cookie http.Cookie)     //设置响应的cookie
	RemoveCookie(cookieName string)    //删除cookie
}

type Responser interface {
	SetHeader(key, value string)          //设置响应头
	AddHeader(key, value string)          //给响应头添加值
	SetContentType(contentType string)    //给响应头设置contenttype
	SetStatusCode(code int)               //设置响应头的状态码
	Header() http.Header                  //返回响应头信息
	HeaderValue(key string) string        //返回某一字段的响应头信息
	Redirect(code int, targetUrl string)  //重定向
	Send(data []byte) (size int)          //给客户端发送响应,返回发送的字节长度
	SendJson(data interface{}) (size int) //响应json格式的数据,返回发送的字节长度
	SetCookie(cookie http.Cookie)         //设置cookie
	//Hijack() (*HijackUp, error)  //将response升级为hijack，升级后，Send,SendJson等的相应方法也会调用hijack相关的
}

var (
	_ Responser = &Response{}
	_ Contexter = &Context{}
	_ Responser = &HijackUp{}
)

type Context struct {
	request             *Request
	response            *Response
	handlerList         *list.List
	currentHandler      *list.Element
	ctx                 ctxt.Context //标准包的context
	cancel              ctxt.CancelFunc
	isHijack            bool
	hijacker            *HijackUp
	isWriteTimeout      bool
	isReadyWriteTimeout bool //通知作用，通知该处理即将进入超时状态并处理超时逻辑
	isEnd               bool
	done                chan int
	ended chan int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		request:     NewRequest(r),
		response:    NewResponse(w),
		handlerList: list.New(),
	}
}

func (c *Context) reset(w *Response, r *Request) *Context {
	c.request = r
	c.response = w
	c.handlerList = list.New()
	c.isHijack = false
	c.isWriteTimeout = false
	c.hijacker = nil
	c.isEnd = false
	c.isReadyWriteTimeout = false
	c.done = make(chan int)
	c.ended = make(chan int, 1)
	//c.ctx, c.cancel = ctxt.WithTimeout(ctxt.Background(), writerTimeout)
	go c.timeout()
	return c
}

//释放context里面的信息
func (c *Context) release() {
	c.response = nil
	c.request = nil
	c.currentHandler = nil
	c.handlerList = list.New()
	c.ctx = ctxt.Background()
	c.isHijack = false
	c.hijacker = nil
	c.isReadyWriteTimeout = false
	c.isWriteTimeout = false
	c.isEnd = false
	c.done = nil
	c.ended = nil
	c.hijacker = nil
}

//请求超时
func (c *Context) timeout() {
	urlPath := c.request.Path()
	method := c.request.Method()
	select {
	case <-c.done:
		//发送响应数据之前如果未超时就不算超时，发送数据的过程不算入时间
		return
	case <-time.After(writerTimeout):
		//响应超时了
		c.isReadyWriteTimeout = true
		HttpError(c, method+" "+urlPath+" response timeout", 504)
		c.isWriteTimeout = true
	}
}

//处理下一个中间件
func (c *Context) Next(err ...error) error {
	if c.currentHandler == nil {
		//未知错误
		return errors.New("unknown error of nil handler")
	}
	next := c.currentHandler.Next()
	if next == nil {
		//处理器已处理完毕，或者一些位置错误
		return errors.New("all handler was done")
	}

	if len(err) > 0 && err[0] != nil {
		HttpError(c, err[0].Error(), http.StatusInternalServerError)
		return nil
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

//终止http操作，一般是终止响应操作
func (c *Context) CancelResponse() {
	c.cancel()
}

//获取request的query  map[string][]string
func (c *Context) Query() url.Values {
	return c.request.Query()
}

//根据key查询query
func (c *Context) QueryString(key string) string {
	return c.request.QueryString(key)
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
func (c *Context) FormFile(key string) (*riderServer.UploadFile, error) {
	return c.request.FormFile(key)
}

//根据客户端传过来的字段名返回文件列表
func (c *Context) FormFiles(key string) ([]*riderServer.UploadFile, error) {
	return c.request.FormFiles(key)
}

//保存客服端传过来的文件（multipart/form-data）
func (c *Context) StoreFormFile(file *riderServer.UploadFile, fileName string) (int64, error) {
	return c.request.StoreFormFile(file, fileName)
}

//保存客户端传过来的文件列表（multipart/form-data）
func (c *Context) StoreFormFiles(files []*riderServer.UploadFile, path string) ([]string, error) {
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

//获取requestID
func (c *Context) RequestID() string {
	return c.request.requestID
}

//获取请求头
func (c *Context) Header() map[string][]string {
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
func (c *Context) RemoteIP() string {
	return c.request.RemoteIP()
}

//RemoteAddr to an "IP:port" address
func (c *Context) FullRemoteIP() string {
	return c.request.FullRemoteIP()
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
func (c *Context) SetContentType(contentType string) {
	if c.isHijack {
		c.hijacker.SetHeader("Content-Type", contentType)
	} else {
		c.response.SetContentType(contentType)
	}
}

//设置响应的状态码
func (c *Context) SetStatusCode(code int) {
	if c.isHijack {
		c.hijacker.SetStatusCode(code)
	} else {
		c.response.SetStatusCode(code)
	}
}

//获取响应头信息
func (c *Context) SendHeader() http.Header {
	if c.isHijack {
		return c.hijacker.Header()
	} else {
		return c.response.Header()
	}
}

//获取响应头的某一字段的值
func (c *Context) SendHeaderValue(key string) string {
	if c.isHijack {
		return c.hijacker.HeaderValue(key)
	} else {
		return c.response.HeaderValue(key)
	}
}

//设置cookies
func (c *Context) SendCookie(cookie http.Cookie) {
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
	c.SendCookie(cookie)
}

/*
	if c.isWriteTimeout {
		return
	}
	如果超时了调用SendJson(), Redirect(), Send()不会再像客户端发送消息
*/

//给客户端发送响应
func (c *Context) Send(data []byte) {
	if c.isWriteTimeout {
		return
	}
	c.End()
	if c.isHijack {
		c.hijacker.Send(data)
	} else {
		c.response.Send(data)
	}
}

//发送json格式数据给客户端
func (c *Context) SendJson(data interface{}) {
	if c.isWriteTimeout {
		return
	}
	c.End()
	if c.isHijack {
		c.hijacker.SendJson(data)
	} else {
		c.response.SendJson(data)
	}
}

//重定向
func (c *Context) Redirect(code int, targetUrl string) {
	if c.isWriteTimeout {
		return
	}
	c.End()
	if c.isHijack {
		c.hijacker.Redirect(code, targetUrl)
	} else {
		c.response.Redirect(code, targetUrl)
	}
}

//通知响应结束
func (c *Context) End() {
	if c.isEnd {
		return
	}
	c.isEnd = true
	if !(c.isReadyWriteTimeout) {
		//当状态进入超时时，不会给done传递信号
		c.done <- 0
	}
	close(c.done)
	c.ended <- 0
	close(c.ended)

	//c.CancelResponse()
}

//hijack相关实现

//升级responsewrite为hijack
func (c *Context) Hijack() (*HijackUp, error) {
	originHeader := c.SendHeader()
	hj, ok := c.response.writer.(http.Hijacker)
	if !ok {
		return nil, errors.New("服务不支持升级hijack")
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		return nil, errors.New("Hijack error: " + err.Error())
	}
	c.isHijack = true
	c.hijacker = &HijackUp{conn: conn, bufrw: bufrw, header: originHeader}
	return c.hijacker, nil
}
