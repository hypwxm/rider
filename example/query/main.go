package main

import (
	"rider"
)

func main() {
	app := rider.New()
	app.GET("/query", func (c rider.Context) {
			//c.Send([]byte(c.QueryString("a")))
			c.SendJson(200, c.Query())
	})
	app.Listen(":5006")
}
