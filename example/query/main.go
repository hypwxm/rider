package main

import (
	"rider"
)

func main() {
	app := rider3.New()
	app.GET("/query", &rider3.Router{
		Handler: func (c *rider3.Context) {
			//c.Send([]byte(c.QueryString("a")))
			c.SendJson(c.Query())
		},
	})
	app.Listen(":5006")
}
