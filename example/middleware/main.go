package main

import (
	"rider"
	"fmt"
	"rider/example/middleware/subMiddle"
)

func main() {
	app := rider.New()
	app.AddMiddleware(
		func(c *rider.Context) {
			fmt.Println("1")
			c.Next()
		},
	)
	app.GET("/middle", &rider.Router{
		Handler: func(c *rider.Context) {
			fmt.Println("2")
			c.Send([]byte("ok"))
		},
	})
	app.GET("/sub", subMiddle.Router())
	app.Listen(":5009")
}
