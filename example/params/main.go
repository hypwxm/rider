package main

import (
	"fmt"

	"github.com/hypwxm/rider"
)

func main() {
	//目前不支持params和pathParams混写
	app := rider.New()
	app.GET("/:id/id2/:id3", func(c rider.Context) {
		fmt.Println(c.Params())
		fmt.Println(c.Param("id"))
		fmt.Println(c.Param("id2"))
		c.Send(200, []byte("ok1"))
	})

	app.GET("/id/id2/(.*)", func(c rider.Context) {
		fmt.Println(c.Params())
		fmt.Println(c.PathParams())
		c.Send(200, []byte("ok2"))
	})

	//pathParams解析正则匹配到的数据
	app.GET("/reg/([a-z]{1}([a-z]{1}))/(.*)", func(c rider.Context) {
		fmt.Println(c.PathParams())
		c.Send(200, []byte("ok3"))
	})
	app.Listen(":5001")
}
