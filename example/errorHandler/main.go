package main

import (
	"rider"
	"errors"
	"fmt"
)

func main() {
	app := rider.New()
	app.GET("/error", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Next(errors.New("错误处理自定义"))
		},
	})
	app.Error(func(c *rider.Context, err string, code int) {
		fmt.Println(err)
		c.Send([]byte(err))
	})
	app.Listen(":5005")
}
