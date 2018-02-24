package subMiddle

import (
	"fmt"
	"rider"
)

func RouterSub() *rider2.Router {
	subrouter := rider2.NewRouter()
	subrouter.GET("/midd", func(c rider2.Context) {
		fmt.Println("mide in mid in mid2")
	})

	return subrouter
}
