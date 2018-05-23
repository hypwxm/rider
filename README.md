# rider
A lightweight http framework, custom infinite pole routing

## 下载
* go get rider

### 1: 创建服务
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
## 5: render模版
* [Render](https://github.com/hypwxm/rider/tree/master/example/render)
## 6: 设置静态文件
* [static](https://github.com/hypwxm/rider/tree/master/example/setstatic)
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
## 7: 接口

```go
	//获取Response响应体内容
	Response() *Response
```
```go
	//获取请求体部分
	Request() *Request
```
* [logger](https://github.com/hypwxm/rider/tree/master/example/logger)
```go
	//获取app初始化时注册的logger服务，用于生成日志
	Logger() *logger.LogQueue
```
* 见下文
```go
	//###获取请求相关的数据，
	//给context传递变量，该变量在整个请求的传递中一直有效（中间件传递）
	SetLocals(key string, value interface{})

	//通过SetLocals设置的值可以在整个响应处理环节通过GetLocals获取
	GetLocals(key string) interface{}
```
* [Query](https://github.com/hypwxm/rider/tree/master/example/query)
* 注. Query只获取url里面的参数。
```go
	//只获取请求url中的查询字符串querystring的map[]string值
	Query() url.Values

	//根据字段名直接查询querystring某个字段名对应的值
	QueryString(key string) string
```
* Body(). 只获取请求体的内容，包括form格式个multipart格式，（不包含上传文件）

```go
	//只获取请求体内的请求参数，
	Body() url.Values

	//根据字段名直接查询"请求体"某中个字段名对应的值
	BodyValue(key string) string
```
* [Params](https://github.com/hypwxm/rider/tree/master/example/params)
```go
	//获取请求路由 /:id/:xx中值，路由参数
	Params() map[string]string

	//获取请求路由中某字段的值
	Param(key string) string
```
```go
	//注册的请求路径中存在正则匹配规则，取得的参数和取正则sub参数一样（详见param例子）
	PathParams() []string //通过正则匹配后得到的路径上的一些参数
```
* [uploadFile](https://github.com/hypwxm/rider/tree/master/example/uploadFile)
```go
	//当请求头的content-type为multipart/form-data时，获取请求中key对应的文件信息(多个文件时，只会获取第一个)
	FormFile(key string) (*UploadFile, error)

	//当请求头的content-type为multipart/form-data时，获取请求中key对应的文件列表
	FormFiles(key string) ([]*UploadFile, error)

	//将file保存，fileName指定完整的路径和名称（先调用FormFile或者FormFiles将返回的file传入即可）
	StoreFormFile(file *UploadFile, fileName string) (int64, error)

	//先通过FormFiles获取文件列表，指定path目录，存储文件的文件夹。文件名将会用随机字符串加文件的后缀名（file.GetFileExt()）
	StoreFormFiles(files []*UploadFile, path string) ([]string, error)
```
* [Header](https://github.com/hypwxm/rider/tree/master/example/header)
```go
	//获取请求头信息
	Header() http.Header

	//根据key获取请求头某一字段的值
	HeaderValue(key string) string
```
```go
	//获取请求头的完整url
	URL() string
```
```go
	//获取请求头的path
	Path() string
```
```go
	//获取响应头的HTTP方法
	Method() string
```
```go
	//获取请求来源的IP地址
	ClientIP() string
```
```go
	//获取分配的请求id
	RequestID() string
```
```go
	//判断请求是否为ajax
	IsAjax() bool
```
```go
	//获取请求的cookies
	Cookies() []*http.Cookie

	//获取请求体中某一字段的cookie值
	CookieValue(key string) (string, error)
```
```go
	//##响应相关
	//获取响应的状态码
	Status() int
```
* [Hijack](https://github.com/hypwxm/rider/tree/master/example/hijack)
```go
	//将http请求升级为hijack，hijack的信息保存在HijackUp中
	Hijack() (*HijackUp, error)
```
* [SendFile](https://github.com/hypwxm/rider/tree/master/example/sendFile)
```go
	//给客户端发送文件（可用户静态文件处理）
	SendFile(path string) error
```
* [download](https://github.com/hypwxm/rider/tree/master/example/download)
```go
	//下载，fileName为完整路径,name为下载时指定的下载名称，传""使用文件本身名称，typ指定attachment还是inline的方式下载，默认为attchment
	Download(fileName string, name string, typ string) error
```
```go
	//获取完整的响应头信息
	WHeader() http.Header

	//获取响应头的某一字段值
	WHeaderValue(key string) string
```
* [Cookie](https://github.com/hypwxm/rider/tree/master/example/setCookie)

```go
	//设置cookie
	SetCookie(cookie http.Cookie)

	//根据cookie名删除cookie
	DeleteCookie(cookieName string)
```
* [Jwt](https://github.com/hypwxm/rider/tree/master/example/jwt)
```go
	//获取系统默认注册的jwt服务（"github.com/dgrijalva/jwt-go"），具体使用例子参考example/jwt
	Jwt() *riderJwter

	//注入jwt（使用前必须先注入服务）
	setJwt(jwter *riderJwter)
```
```go
	//设置响应头
	SetHeader(key, value string)

	//给响应头添加值
	AddHeader(key, value string)
```
```go
	//给响应头设置contenttype
	SetCType(contentType string)
```
* [Redirect](https://github.com/hypwxm/rider/tree/master/example/redirect)
```go
	//重定向
	Redirect(code int, targetUrl string)
```
* [Send & SendJson](https://github.com/hypwxm/rider/tree/master/example/send)
```go
	//给客户端发送数据
	Send(code int, d []byte) (int, error)

	//给客户端发送json格式的数据
	SendJson(code int, i interface{}) (int, error)
```
* [Render](https://github.com/hypwxm/rider/tree/master/example/render)
```go
	//负责模板渲染 ，只要实现了BaseRender，注册app服务是直接修改tplsRender的值
	Render(tplName string, data interface{})
```
## 8: 上下文变量
* [locals](https://github.com/hypwxm/rider/tree/master/example/locals)
*  main.go
```go
app := rider.New()
app.USE(
	func(c rider.Context) {
		c.SetLocals("locals", "this is the first locals")
		c.SetLocals("locals2", "this is the second locals")
		c.Next()
	},
)
app.Kid("/", router.Router())
app.Listen(":5003")
```
*  child.go
```go
func Router() *rider.Router {
	router := rider.NewRouter()
	router.GET("/", func (c rider.Context) {
		fmt.Println(c.GetLocals("locals"))
		fmt.Println(c.GetLocals("locals2"))
		c.Send(200, []byte("ok"))
	})
	router.GET("/xx", func (c rider.Context) {
		fmt.Println(c.GetLocals("locals"))
		fmt.Println(c.GetLocals("locals2"))
		c.Send(200, []byte("ok2"))
	})
	return router
}
```
## 9: 下载模块
* [download](https://github.com/hypwxm/rider/tree/master/example/download)
```go c.Download(filename, name, type) ```
* @params
* 1: filename: 文件所在路径加全名
* 2: name: 指定下载时文件的名称，不指定默认取路径中的名称
* 3: type: 下载文件的方式，attachment和inline。
```go
app.GET("/download", func(c rider.Context) {
	c.Download(filename, name, type)
})
```
## 10: jwt模块（可取代cookie），比cookie要安全
* [jwt](https://github.com/hypwxm/rider/tree/master/example/jwt)
* 作为中间件引入
```go
app := rider.New()
app.USE(rider.RiderJwt("rider", time.Hour))
app.GET("/token", func(c rider.Context) {
	token, _ := c.Jwt().Set("test", " test")
	c.Send(200, []byte(token))
})
app.GET("/tokenparse", func(c rider.Context) {
	c.Logger().INFO(c.CookieValue("token"))
	c.Jwt().Delete("test")
	c.Jwt().DeleteAll()
	c.Jwt().Set("a", "b")
	fmt.Println(c.Jwt().Claims())
	fmt.Println(c.Jwt().ClaimsValue("a"))
})
app.Listen(":5002")
```
## 11: logger日志模块
* [logger](https://github.com/hypwxm/rider/tree/master/example/logger)
* 注册服务后，日志模块会一同注册。
* 通过smtp服务可以注册邮箱日志
* 调用```go app.Logger(8) ``` 修改日志等级
* 日志等级分为
 1): fatalLevel      uint8 = iota 
 * 日志打印的同时服务也会退出
 ```go 
 c.Logger().FATAL("")
 ```
 2): panicLevel
 ```go 
 c.Logger().PANIC("")
 ```
 3): errorLevel
 ```go 
 c.Logger().ERROR("")
 ```
 4): warningLevel
 ```go 
 c.Logger().WARNING("")
 ```
 5): infoLevel
 ```go 
 c.Logger().INFO("")
 ```
 6): consoleLevel
 ```go 
 c.Logger().CONSOLE("")
 ```
 7): debugLevel
 ```go 
 c.Logger().DEBUG("")
 ```
 
```go
app := rider.New()
rlog := app.Logger(8)
```
## 12: smtp邮箱模块
* [smtp](https://github.com/hypwxm/rider/tree/master/example/smtp)
