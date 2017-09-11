package riderRouter


//实现这个接口，就可以将其作为rider实例和Router实例的第二个参数
type RiderHandler interface {
	ServeHTTP(context Context)
}

type HandlerFunc func(context Context)

func (h HandlerFunc) ServeHTTP(context Context) {
	h(context)
}