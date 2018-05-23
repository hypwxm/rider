package main

import (
	"fmt"
	"rider"
)

func main() {
	app := rider.New()
	app.GET("/error", func(c rider.Context) {
		c.Next(rider.NError{500, "错误处理自定义"})
	})
	app.Error(func(c rider.Context, err string, code int) {
		fmt.Println(err)
		c.Send(200, []byte(err))
	})
	app.Listen(":5005")
}
