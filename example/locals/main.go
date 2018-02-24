package main

import (
	"rider"
	"rider/example/locals/router"
)

func main() {
	app := rider2.New()
	app.Logger(8)
	app.USE(
		func(c rider2.Context) {
			c.SetLocals("locals", "this is the first locals")
			c.SetLocals("locals2", "this is the second locals")
			c.Next()
		},
	)
	app.Kid("/", router.Router())
	app.Listen(":5003")
}
