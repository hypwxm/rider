package subMiddle

import (
	"rider"
	"fmt"
)

func RouterSub3() *rider.Router {
	subrouter := rider.NewRouter()
	subrouter.USE(func(c rider.Context) {
		fmt.Println("sub3 sub3")
		c.Next()
	})
	subrouter.GET("/midd", func(c rider.Context) {
		fmt.Println("mide in mid in mid2")
	})

	return subrouter
}
