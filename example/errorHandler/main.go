package main

import (
	"errors"
	"fmt"
	"rider"
)

func main() {
	app := rider2.New()
	app.GET("/error", func(c rider2.Context) {
		c.Next(errors.New("错误处理自定义"))
	})
	app.Error(func(c rider2.Context, err string, code int) {
		fmt.Println(err)
		c.Send(200, []byte(err))
	})
	app.Listen(":5005")
}
