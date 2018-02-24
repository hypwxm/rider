package subMiddle

import (
	"fmt"
	"rider"
)

func RouterSub3() *rider2.Router {
	subrouter := rider2.NewRouter()
	subrouter.USE(func(c rider2.Context) {
		fmt.Println("sub3 sub3")
		c.Next()
	})
	subrouter.GET("/midd", func(c rider2.Context) {
		fmt.Println("mide in mid in mid2")
	})

	return subrouter
}
