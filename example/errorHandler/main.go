package main

import (
	"rider"
	"errors"
	"fmt"
)

func main() {
	app := rider3.New()
	app.GET("/error", &rider3.Router{
		Handler: func(c *rider3.Context) {
			c.Next(errors.New("错误处理自定义"))
		},
	})
	app.Error(func(c *rider3.Context, err string, code int) {
		fmt.Println(err)
		c.Send([]byte(err))
	})
	app.Listen(":5005")
}
