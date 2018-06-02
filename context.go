package rider

import (
	"container/list"
	// ctxt "context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/hypwxm/rider/logger"
)

type NError struct {
	Status int
	Error  string
}

type Context interface {
	//设置context的任务函数处理队列（包含中间件）
	handlerQueue(list *list.List)

	//实现中间件的链表处理
	Next(err ...NError) error

	//设置链表中当前在处理处理器
	setCurrent(element *list.Element)

	//设置开始处理的第一个中间件
	setStartHandler() *list.Element

	//获取当前处理中的中间件
	getCurrentHandler() (HandlerFunc, error)

	//开始处理中间件，开始响应
	startHandleList() error

	//释放context
	release()

	//初始化context
	reset(w *Response, r *Request, server *HttpServer) *context

	//获取Response响应体内容
	Response() *Response

	//获取请求体部分
	Request() *Request

	//获取app初始化时注册的logger服务，用于生成日志
	Logger() *logger.LogQueue

	//###获取请求相关的数据，
	//给context传递变量，该变量在整个请求的传递中一直有效（中间件传递）
	SetLocals(key string, value interface{})

	//通过SetLocals设置的值可以在整个响应处理环节通过GetLocals获取
	GetLocals(key string) interface{}

	Locals() map[string]interface{} //获取SetLocals设置的所有变量，用户输出给模板

	// 删除locals里面的数据
	DeleteLocals(key string)

	// 删除locals的所有信息
	DeleteAllLocals()

	//只获取请求url中的查询字符串querystring的map[]string值
	Query() url.Values

	//根据字段名直接查询querystring某个字段名对应的值
	QueryString(key string) string

	// 当query不存在某一字段的值时，返回一个默认的值
	QueryDefault(key string, def string) string

	//只获取请求体内的请求参数，
	Body() url.Values

	//根据字段名直接查询"请求体"某中个字段名对应的值
	BodyValue(key string) string

	// 当body不存在某一字段的值时，返回一个默认的值
	BodyDefault(key string, def string) string

	// 返回请求的参数（包括query和body的）
	Form() url.Values

	//根据字段名直接查询"请求参数"某中个字段名对应的值
	FormValue(key string) string

	// 当请求参数不存在某一字段的值时，返回一个默认的值
	FormDefault(key string, def string) string

	//获取请求路由 /:id/:xx中值，路由参数
	Params() map[string]string

	//获取请求路由中某字段的值
	Param(key string) string

	//注册的请求路径中存在正则匹配规则，取得的参数和取正则sub参数一样（详见param例子）
	PathParams() []string //通过正则匹配后得到的路径上的一些参数

	//当请求头的content-type为multipart/form-data时，获取请求中key对应的文件信息(多个文件时，只会获取第一个)
	FormFile(key string) (*UploadFile, error)

	//当请求头的content-type为multipart/form-data时，获取请求中key对应的文件列表
	FormFiles(key string) ([]*UploadFile, error)

	//将file保存，fileName指定完整的路径和名称（先调用FormFile或者FormFiles将返回的file传入即可）
	StoreFormFile(file *UploadFile, fileName string) (int64, error)

	//先通过FormFiles获取文件列表，指定path目录，存储文件的文件夹。文件名将会用随机字符串加文件的后缀名（file.GetFileExt()）
	StoreFormFiles(files []*UploadFile, path string) ([]string, error)

	//获取请求头信息
	Header() http.Header

	//根据key获取请求头某一字段的值
	HeaderValue(key string) string

	//获取请求头的完整url
	URL() string

	//获取请求头的path
	Path() string

	//获取响应头的HTTP方法
	Method() string

	//获取请求来源的IP地址
	ClientIP() string

	//获取分配的请求id
	RequestID() string

	//判断请求是否为ajax
	IsAjax() bool

	//获取请求的cookies
	Cookies() []*http.Cookie

	//获取请求体中某一字段的cookie值
	CookieValue(key string) (string, error)

	//##响应相关
	//获取响应的状态码
	Status() int

	//将http请求升级为hijack，hijack的信息保存在HijackUp中
	Hijack() (*HijackUp, error)

	//给客户端发送文件（可用户静态文件处理）
	SendFile(path string) error

	//下载，fileName为完整路径,name为下载时指定的下载名称，传""使用文件本身名称，typ指定attachment还是inline的方式下载，默认为attchment
	Download(fileName string, name string, typ string) error

	//获取完整的响应头信息
	WHeader() http.Header

	//获取响应头的某一字段值
	WHeaderValue(key string) string

	//设置cookie
	SetCookie(cookie http.Cookie)

	//根据cookie名删除cookie
	DeleteCookie(cookieName string)

	//获取系统默认注册的jwt服务（"github.com/dgrijalva/jwt-go"），具体使用例子参考example/jwt
	Jwt() *riderJwter

	//注入jwt（使用前必须先注入服务）
	setJwt(jwter *riderJwter)

	//设置响应头
	SetHeader(key, value string)

	//给响应头添加值
	AddHeader(key, value string)

	//给响应头设置contenttype
	SetCType(contentType string)

	//重定向
	Redirect(code int, targetUrl string)

	//给客户端发送数据
	Send(code int, d []byte) (int, error)

	//发送字符串
	SendString(code int, s string) (int, error)

	//给客户端发送json格式的数据
	SendJson(code int, i interface{}) (int, error)

	//负责模板渲染 ，只要实现了BaseRender，注册app服务是直接修改tplsRender的值
	Render(tplName string, data interface{})

	// 获取httpserver实例
	getHttpServer() *HttpServer

	// 获取请求来源host
	Host() string
}

var (
	_ Context = &context{}
)

// "##"表示从pool取得context时必须初始化的值
type context struct {
	request        *Request      //##
	response       *Response     //##
	handlerList    *list.List    //##
	currentHandler *list.Element //##
	// ctx            ctxt.Context           //##标准包的context
	// cancel         ctxt.CancelFunc        //##
	isHijack  bool                   //##
	hijacker  *HijackUp              //##
	isEnd     bool                   //##
	done      chan int               //##
	server    *HttpServer            //整个服务引用的server都是同一个
	ended     chan int               //##
	committed bool                   //##表示状态码是否已经发送（writeHeader有无调用）
	jwt       *riderJwter            //用于存储jwt
	locals    map[string]interface{} //用户存储locals的变量，用户输出给模板时调用，对该变量的所有操作都未加锁，如需多协程读写，请自行加锁
	query     url.Values             //存放请求查询参数
	body      url.Values             //存放请求体内的字段（不包含get查询参数字段）
	form      url.Values             //存放请求参数（包含get，post，put）
}

func newContext(w http.ResponseWriter, r *http.Request, server *HttpServer) Context {
	//从pool取得context
	ctx := basePool.context.Get().(Context)
	//从pool取得requset（注意初始化）
	request := NewRequest(r)
	//从pool取得response（注意初始化）
	response := NewResponse(w, server)
	//（初始化context）
	ctx.reset(response, request, server)
	return ctx
}

func releaseContext(c Context) {
	//响应结束，释放request，response，context回poo，方便其他请求取得（释放时记得参数初始化）
	c.Response().release()
	basePool.response.Put(c.Response())
	c.Request().release()
	basePool.request.Put(c.Request())
	c.release()
	basePool.context.Put(c)
}

func (c *context) reset(w *Response, r *Request, server *HttpServer) *context {
	c.request = r
	c.response = w
	c.currentHandler = nil
	c.handlerList = list.New()
	//c.ctx = ctxt.Background()
	//c.ctx = c.request.request.Context()
	c.isHijack = false
	c.hijacker = nil
	//c.isEnd = false
	c.done = make(chan int)
	c.ended = make(chan int, 1)
	c.server = server
	c.committed = false
	c.jwt = nil
	c.locals = make(map[string]interface{})
	//c.ctx, c.cancel = ctxt.WithTimeout(ctxt.Background(), writerTimeout)
	//go c.timeout()
	c.query = c.request.query()
	c.body = c.request.body()
	c.form = c.request.form()
	return c
}

//释放context里面的信息
func (c *context) release() {
	c.response = nil
	c.request = nil
	c.currentHandler = nil
	c.handlerList = list.New()
	// c.ctx = nil
	c.isHijack = false
	c.hijacker = nil
	c.isEnd = false
	c.done = nil
	//c.ended = nil
	c.hijacker = nil
	c.jwt = nil
	c.committed = false
	c.locals = nil
	c.query = nil
	c.body = nil
	c.form = nil
}

//处理下一个中间件
func (c *context) Next(err ...NError) error {
	if c.currentHandler == nil {
		//未知错误
		return errors.New("unknown error of nil handler")
	}

	if len(err) > 0 {
		status := http.StatusInternalServerError
		if err[0].Status != 0 {
			status = err[0].Status
		}
		HttpError(c, err[0].Error, status)
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

func (c *context) Request() *Request {
	return c.request
}

func (c *context) Jwt() *riderJwter {
	return c.jwt
}

func (c *context) setJwt(jwter *riderJwter) {
	c.jwt = jwter
}

func (c *context) Logger() *logger.LogQueue {
	return c.server.logger
}

func (c *context) handlerQueue(l *list.List) {
	c.handlerList = l
}

//设置链表中当前在处理处理器
func (c *context) setCurrent(e *list.Element) {
	c.currentHandler = e
}

//设置链表中开始处理是的起始handler
func (c *context) setStartHandler() *list.Element {
	e := c.handlerList.Front()
	c.currentHandler = e
	return e
}

//获取当前正在处理中的handler
func (c *context) getCurrentHandler() (HandlerFunc, error) {
	if c.currentHandler != nil {
		return c.currentHandler.Value.(HandlerFunc), nil
	}
	return nil, errors.New("处理函数为nil")
}

//开始处理链表的中的处理器
func (c *context) startHandleList() error {
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
func (c *context) SetLocals(key string, val interface{}) {
	// c.ctx = ctxt.WithValue(c.ctx, key, val)
	c.locals[key] = val
}

//从context中获取临时信息
func (c *context) GetLocals(key string) interface{} {
	// return c.ctx.Value(key)
	return c.locals[key]
}

// 删除context的临时信息
func (c *context) DeleteLocals(key string) {
	if _, ok := c.locals[key]; ok {
		delete(c.locals, key)
	}
}

// 删除context的所有临时变量
func (c *context) DeleteAllLocals() {
	c.locals = make(map[string]interface{})
}

//获取locals数据
func (c *context) Locals() map[string]interface{} {
	return c.locals
}

//获取request的query  map[string][]string
func (c *context) Query() url.Values {
	return c.query
}

//根据key查询query
func (c *context) QueryString(key string) string {
	return c.Query().Get(key)
}

// QueryString返回空时提供提供一个默认值
func (c *context) QueryDefault(key string, def string) string {
	if c.QueryString(key) == "" {
		return def
	} else {
		return c.QueryString(key)
	}
}

//获取request请求体  map[string][]string
func (c *context) Body() url.Values {
	return c.body
}

//根据key查询body里面的某一字段的第一个值
func (c *context) BodyValue(key string) string {
	return c.body.Get(key)
}

// QueryString返回空时提供提供一个默认值
func (c *context) BodyDefault(key string, def string) string {
	if c.BodyValue(key) == "" {
		return def
	} else {
		return c.BodyValue(key)
	}
}

//根据key查询请求参数（包含get，post，put的所有字段）
func (c *context) Form() url.Values {
	return c.form
}

//根据key查询请求参数里面的某一字段的第一个值（包含get，post，put的所有字段）
func (c *context) FormValue(key string) string {
	return c.form.Get(key)
}

// formstring返回空时提供提供一个默认值
func (c *context) FormDefault(key string, def string) string {
	if c.BodyValue(key) == "" {
		return def
	} else {
		return c.BodyValue(key)
	}
}

//根据客户端传过来的字段名会返回第一个文件
func (c *context) FormFile(key string) (*UploadFile, error) {
	return c.request.FormFile(key)
}

//根据客户端传过来的字段名返回文件列表
func (c *context) FormFiles(key string) ([]*UploadFile, error) {
	return c.request.FormFiles(key)
}

//保存客服端传过来的文件（multipart/form-data）
func (c *context) StoreFormFile(file *UploadFile, fileName string) (int64, error) {
	return c.request.StoreFormFile(file, fileName)
}

//保存客户端传过来的文件列表（multipart/form-data）
func (c *context) StoreFormFiles(files []*UploadFile, path string) ([]string, error) {
	return c.request.StoreFormFiles(files, path)
}

//获取path匹配的所有参数/:id/:xx
func (c *context) Params() map[string]string {
	return c.request.Params()
}

//获取path匹配的参数/:id/:xx
func (c *context) Param(key string) string {
	return c.request.Param(key)
}

//获取正则匹配路径后的一些参数
func (c *context) PathParams() []string {
	return c.request.pathParams
}

//获取requestID
func (c *context) RequestID() string {
	return c.request.requestID
}

//获取请求头
func (c *context) Header() http.Header {
	return c.request.Header()
}

//获取请求头的某一字段的值
func (c *context) HeaderValue(key string) string {
	return c.request.HeaderValue(key)
}

//获取请求的URL
func (c *context) URL() string {
	return c.request.Url()
}

//获取请求的PATH
func (c *context) Path() string {
	return c.request.Path()
}

//获取请求的HTTP方法
func (c *context) Method() string {
	return c.request.Method()
}

//RemoteAddr to an "IP" address
func (c *context) ClientIP() string {
	return c.request.ClientIP()
}

func (c *context) Response() *Response {
	return c.response
}

//判断请求是否为ajax
func (c *context) IsAjax() bool {
	return c.request.IsAJAX()
}

//获取请求的cookies
func (c *context) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

//获取某一字段的cookie
func (c *context) CookieValue(key string) (string, error) {
	return c.request.CookieValue(key)
}

//设置响应头
func (c *context) SetHeader(key, value string) {
	if c.isHijack {
		c.hijacker.SetHeader(key, value)
	} else {
		c.response.SetHeader(key, value)
	}
}

//添加响应头
func (c *context) AddHeader(key, value string) {
	if c.isHijack {
		c.hijacker.AddHeader(key, value)
	} else {
		c.response.AddHeader(key, value)
	}
}

//设置响应的contenttype
func (c *context) SetCType(contentType string) {
	if c.isHijack {
		c.hijacker.SetHeader("Content-Type", contentType)
	} else {
		c.response.SetCType(contentType)
	}
}

//获取响应头信息
func (c *context) WHeader() http.Header {
	if c.isHijack {
		return c.hijacker.Header()
	} else {
		return c.response.Header()
	}
}

//获取响应头的某一字段的值
func (c *context) WHeaderValue(key string) string {
	if c.isHijack {
		return c.hijacker.HeaderValue(key)
	} else {
		return c.response.HeaderValue(key)
	}
}

//设置cookies
func (c *context) SetCookie(cookie http.Cookie) {
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
func (c *context) DeleteCookie(cookieName string) {
	cookie := http.Cookie{Name: cookieName, MaxAge: -1, Path: "/"}
	c.SetCookie(cookie)
}

//
func (c *context) CloseNotify() <-chan bool {
	return c.response.CloseNotify()
}

//获取响应状态码
func (c *context) Status() int {
	if c.isHijack {
		return c.hijacker.status
	} else {
		return c.response.status
	}
}

func (c *context) writeHeader(code int) {
	if c.isHijack {
		c.hijacker.WriteHeader(code)
	} else {
		c.response.WriteHeader(code)
	}
}

//给客户端发送响应
func (c *context) Send(code int, data []byte) (size int, err error) {
	if c.WHeader().Get(HeaderContentType) == "" {
		c.SetCType(http.DetectContentType(data) + ";charset=utf8")
	}

	if code == 200 || code == 304 {
		if setWeakEtag(c, data, c.Request().request) {
			c.writeHeader(304)
			if c.isHijack {
				c.hijacker.Write([]byte(""))
			} else {
				c.response.Write([]byte(""))
			}
			return 0, nil
		}
	}

	c.writeHeader(code)

	if c.isHijack {
		return c.hijacker.Write(data)
	} else {
		return c.response.Write(data)
	}
}

func (c *context) SendString(code int, s string) (int, error) {
	return c.Send(code, []byte(s))
}

//发送json格式数据给客户端
func (c *context) SendJson(code int, data interface{}) (int, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	c.SetCType("application/json; charset=utf8")
	return c.Send(code, jsonData)
}

//重定向
func (c *context) Redirect(code int, targetUrl string) {
	if code < 300 || code > 308 {
		c.server.logger.PANIC("Invalid redirect status code")
		return
	}
	if code == 301 {
		if c.Method() != http.MethodGet {
			code = 307
		}
	}
	c.SetHeader("Content-Length", "0")
	c.response.Redirect(code, targetUrl)
}

//hijack相关实现

//升级responsewrite为hijack
func (c *context) Hijack() (*HijackUp, error) {
	var originHeader http.Header = c.WHeader()
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
func (c *context) Render(tplName string, data interface{}) {
	var err error
	c.SetCType("text/html;charset=utf8")
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
func (c *context) SendFile(path string) error {
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

	chunk, err := ioutil.ReadAll(fp)

	if setWeakEtag(c, chunk, c.request.request) {
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
func (c *context) Download(fileName string, name string, typ string) error {
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

// 获取服务实例
func (c *context) getHttpServer() *HttpServer {
	return c.server
}

// 获取请求来源host
func (c *context) Host() string {
	return c.request.Req().Host
}
