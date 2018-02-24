package main

import (
	"fmt"
	"rider"
	"rider/example/middleware/subMiddle"
)

func mid() rider2.HandlerFunc {
	return func(context rider2.Context) {
		fmt.Println("funcmid")
		context.Next()
	}
}

func main() {
	app := rider2.New()
	app.USE(
		func(c rider2.Context) {
			fmt.Println("1")
			c.Next()
		},
	)
	app.GET("/middle", mid(), func(c rider2.Context) {
		fmt.Println("2")
		c.Send(200, []byte("ok"))
	})

	app.Kid("/sub", mid(), subMiddle.Router())
	app.Kid("/sub", subMiddle.RouterSub())

	app.Listen(":5009")
}
