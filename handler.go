package rider

//实现这个接口，就可以将其作为rider实例和Router实例的第二个参数
type riderHandler interface {
	riderServeHTTP(context *Context)
}

type IsRouterHandler interface {
	BeenHandler()
}

type HandlerFunc func(context Context)

func (h HandlerFunc) BeenHandler() {

}

func (r *Router) BeenHandler() {

}
