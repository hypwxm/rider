package main

import (
	"rider"
	"fmt"
)

func main() {
	app := rider.New()
	app.GET("/", &rider.Router{
		Handler: func(c *rider.Context) {
			//获取请求头信息
			heades := c.Header()
			//获取请求头字段信息，数组的第一个值
			accept := c.HeaderValue("Accept")
			fmt.Println(accept)
			//添加响应头信息
			c.AddHeader("rider", "kick")
			c.SendJson(heades)
		},
	})
	app.Listen(":5004")
}