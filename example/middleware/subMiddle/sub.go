package subMiddle

import (
	"fmt"
	"rider"
)

func RouterSub() *rider.Router {
	subrouter := rider.NewRouter()
	subrouter.GET("/midd", func(c rider.Context) {
		fmt.Println("mide in mid in mid2")
	})

	return subrouter
}
