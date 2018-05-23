package subMiddle

import (
	"fmt"
	"rider"
)

func mid2() rider.HandlerFunc {
	return func(context rider.Context) {
		fmt.Println("funcmid3")
		context.Next()
	}
}

func RouterSub2() *rider.Router {
	subrouter := rider.NewRouter()
	subrouter.GET("/midd2", func(c rider.Context) {
		fmt.Println("mide in mid in mid3")
	})
	subrouter.Kid("/sub3", mid2(), RouterSub3())
	return subrouter
}
