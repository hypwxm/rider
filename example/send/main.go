package main

import (
	"rider"
)

func main() {
	app := rider2.New()
	app.Logger(8)
	app.GET("/", func(c rider2.Context) {
		c.Hijack()
		c.Send(200, []byte("xx"))

	})
	app.Listen(":5001")
}
