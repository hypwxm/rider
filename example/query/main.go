package main

import (
	"rider"
)

func main() {
	app := rider.New()
	app.GET("/query", &rider.Router{
		Handler: func (c *rider.Context) {
			//c.Send([]byte(c.QueryString("a")))
			c.SendJson(c.Query())
		},
	})
	app.Listen(":5006")
}
