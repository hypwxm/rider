package main

import (
	"rider"
)

func main() {
	app := rider.New()
	app.Logger(8)
	app.GET("/", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Hijack()
			c.Send([]byte("xx"))
			c.Send([]byte("xx"))
			c.Send([]byte("xx"))

			c.Send([]byte("xx"))

			c.Send([]byte("xx"))
			c.Send([]byte("xx"))
			c.Send([]byte("xx"))
			c.Send([]byte("xx"))
			c.Send([]byte("xx"))
			c.Send([]byte("xx"))

		},
	})
	app.Listen(":5001")
}
