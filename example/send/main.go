package main

import (
	"github.com/hypwxm/rider"
)

func main() {
	app := rider.New()
	app.Logger(8)
	app.GET("/", func(c rider.Context) {
		c.Hijack()
		c.Send(200, []byte("xx"))

	})
	app.Listen(":5001")
}
