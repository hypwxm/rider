package main

import (
	"rider"
)

func main() {
	app := rider2.New()
	app.Logger(8)
	app.GET("/query", func(c rider2.Context) {
		//c.Send([]byte(c.QueryString("a")))
		//time.Sleep(20e9)
		c.Query().Del("a")
		c.SendJson(200, c.Query())
	})

	app.Graceful(":5006")
}
