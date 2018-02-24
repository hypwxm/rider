package subMiddle

import (
	"fmt"
	"rider"
)

func mid2() rider2.HandlerFunc {
	return func(context rider2.Context) {
		fmt.Println("funcmid3")
		context.Next()
	}
}

func RouterSub2() *rider2.Router {
	subrouter := rider2.NewRouter()
	subrouter.GET("/midd2", func(c rider2.Context) {
		fmt.Println("mide in mid in mid3")
	})
	subrouter.Kid("/sub3", mid2(), RouterSub3())
	return subrouter
}
