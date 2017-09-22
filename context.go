package rider

import (
	"net/http"
	"container/list"
	"errors"
)

type Contexter interface {
	NewContext(w http.ResponseWriter, r *http.Request) //初始化一个Context
	Next()   //实现中间件的链表处理
	SetCurrent(element *list.Element)  //设置链表中当前在处理处理器
	SetStartHandler() *list.Element //设置开始处理的第一个中间件
	GetCurrentHandler() (HandlerFunc, error)  //获取当前处理中的中间件
	StartHandleList() error  //开始处理中间件
	release()  //释放Context
	Reset(w http.ResponseWriter, r *http.Request)  //初始化Context
}

type Context struct {
	Request *Request
	Response *Response
	handlerList *list.List
	currentHandler *list.Element
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request: NewRequest(r),
		Response: NewResponse(w),
		handlerList: list.New(),
	}
}

func (c *Context) reset(w *Response, r *Request) {
	c.Request = r
	c.Response = w
	c.handlerList = list.New()
}

//处理下一个中间件
func (c *Context) Next() {
	if c.currentHandler == nil {
		//未知错误
		c.release()
		return
	}
	next := c.currentHandler.Next()
	if next == nil {
		//处理器已处理完毕，或者一些位置错误
		c.release()
		return
	}
	//先设置current
	c.SetCurrent(next)
	//在执行链表处理函数
	next.Value.(HandlerFunc)(c)
}

//设置链表中当前在处理处理器
func (c *Context) SetCurrent(e *list.Element) {
	c.currentHandler = e
}

//设置链表中开始处理是的起始handler
func (c *Context) SetStartHandler() *list.Element {
	e := c.handlerList.Front()
	c.currentHandler = e
	return e
}

//获取当前正在处理中的handler
func (c *Context) GetCurrentHandler() (HandlerFunc, error) {
	if c.currentHandler != nil {
		return c.currentHandler.Value.(HandlerFunc), nil
	}
	return nil, errors.New("处理函数为nil")
}

//开始处理链表的中的处理器
func (c *Context) StartHandleList() error {
	c.SetCurrent(c.SetStartHandler())
	//开始处理
	handlerFunc, err := c.GetCurrentHandler()
	if err != nil {
		return err
	}
	handlerFunc(c)
	return nil
}

//释放context里面的信息
func (c *Context) release() {
	c.Response = nil
	c.Request = nil
	c.currentHandler = nil
	c.handlerList = list.New()
}