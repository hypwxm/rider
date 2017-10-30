# rider
A lightweight http framework, custom infinite pole routing


## 1: 创建服务
```go
app := rider.New()
```

## 2: 创建路由
```go
app.GET("/path", func(c rider.Context) {
  c.Send(200, []byte("ok"))
})
```
## 3: 添加中间件
```go
app.USE(func(c rider.Context){})
```
```go
app.GET("/path", func(c rider.Context){}, func(c rider.Context){},...)
```
* 只有Kid方法能引入子路由
```go
app.Kid("/path", func(c rider.Context){}, func(c rider.Context){}, func(c rider.Context){}, *rider.Router{})
```
## 4: 路由
* GET/POST/PUT/PATCH/DELETE/OPTIONS
* ANY为任意请求方式的路由
* 通过Kid方式注册子路由，支持无限子路由
## 5: 中间件
* 全局的中间件注册
```go
app := rider.New()
app.USE(func(c rider.Context){})
```
* 路由级的中间件注册
```go
app := rider.New()
app.GET("/path", func(c rider.Context), func(c rider.Context), ...)
```
* 子路由的中间件注册
```go
app := rider.New()
app.Kid("/path", func(c rider.Context), func(c rider.Context), func(c rider.Context), *rider.Router{})
```
## 接口

```go
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

	//只获取请求url中的查询字符串querystring的map[]string值
	Query() url.Values

	//根据字段名直接查询querystring某个字段名对应的值
	QueryString(key string) string

	//只获取请求体内的请求参数，
	Body() url.Values

	//根据字段名直接查询"请求体"某中个字段名对应的值
	BodyValue(key string) string

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

	//给客户端发送json格式的数据
	SendJson(code int, i interface{}) (int, error)

	//负责模板渲染 ，只要实现了BaseRender，注册app服务是直接修改tplsRender的值
	Render(tplName string, data interface{})
  ```
