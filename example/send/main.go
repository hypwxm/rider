package main

import (
	"rider"
)

func main() {
	app := rider3.New()
	app.Logger(8)
	app.GET("/", &rider3.Router{
		Handler: func(c *rider3.Context) {
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
