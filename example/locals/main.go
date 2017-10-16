package main

import (
	"rider"
	"rider/example/locals/router"
)

func main() {
	app := rider.New()
	app.Logger(8)
	app.AddMiddleware(
		func(c *rider.Context) {
			c.SetLocals("locals", "this is the first locals")
			c.SetLocals("locals2", "this is the second locals")
			c.Next()
		},
	)
	app.GET("/", router.Router())
	app.Listen(":5003")
}
