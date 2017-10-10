package main

import (
	"rider"
	"fmt"
)

func main() {
	//目前不支持params和pathParams混写
	app := rider.New()
	app.GET("/:id/id2", &rider.Router{
		Handler: func(c *rider.Context) {
			fmt.Println(c.Params())
			fmt.Println(c.Param("id"))
			fmt.Println(c.Param("id2"))
			c.Send([]byte("ok"))
		},
	})
	//pathParams解析正则匹配到的数据
	app.GET("/reg/([a-z]{1}([a-z]{1}))/(.*)", &rider.Router{
		Handler: func(c *rider.Context) {
			fmt.Println(c.PathParams())
			c.Send([]byte("ok"))
		},
	})
	app.Listen(":5001")
}
