package subMiddle

import (
	"fmt"
	"rider"
)

func mid() rider.HandlerFunc {
	return func(context rider.Context) {
		fmt.Println("funcmid2")
		context.Next()
	}
}

func Router() *rider.Router {
	router := rider.NewRouter()
	router.USE(
		func(c rider.Context) {
			fmt.Println("sub1")
			c.Next()
		},
	)
	router.GET("/subMiddle", func(c rider.Context) {
		fmt.Println("sub2")
		c.Send(200, []byte("ok"))
	})

	router.Kid("/mid", RouterSub())
	router.Kid("/mid", mid(), RouterSub2())

	return router
}
