package subMiddle

import (
	"fmt"
	"rider"
)

func mid() rider2.HandlerFunc {
	return func(context rider2.Context) {
		fmt.Println("funcmid2")
		context.Next()
	}
}

func Router() *rider2.Router {
	router := rider2.NewRouter()
	router.USE(
		func(c rider2.Context) {
			fmt.Println("sub1")
			c.Next()
		},
	)
	router.GET("/subMiddle", func(c rider2.Context) {
		fmt.Println("sub2")
		c.Send(200, []byte("ok"))
	})

	router.Kid("/mid", RouterSub())
	router.Kid("/mid", mid(), RouterSub2())

	return router
}
