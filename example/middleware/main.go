package main

import (
	"rider"
	"fmt"
	"rider/example/middleware/subMiddle"
)

func mid() rider.HandlerFunc {
	return func(context rider.Context) {
		fmt.Println("funcmid")
		context.Next()
	}
}

func main() {
	app := rider.New()
	app.USE(
		func(c rider.Context) {
			fmt.Println("1")
			c.Next()
		},
	)
	app.GET("/middle", func(c rider.Context) {
		fmt.Println("2")
		c.Send(200, []byte("ok"))
	})


	app.Kid("/sub", mid(), subMiddle.Router())
	app.Kid("/sub", subMiddle.RouterSub())

	app.Listen(":5009")
}
