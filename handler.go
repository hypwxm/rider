package rider

//实现这个接口，就可以将其作为rider实例和Router实例的第二个参数
type RiderHandler interface {
	RiderServeHTTP(context *Context)
}

type HandlerFunc func(context *Context)
