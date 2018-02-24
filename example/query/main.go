package main

import (
	"github.com/hypwxm/rider"
)

func main() {
	app := rider.New()
	app.Logger(8)
	app.GET("/query", func(c rider.Context) {
		//c.Send([]byte(c.QueryString("a")))
		//time.Sleep(20e9)
		c.Query().Del("a")
		c.SendJson(200, c.Query())
	})

	app.Graceful(":5006")
}
