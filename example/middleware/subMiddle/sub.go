package subMiddle

import (
	"fmt"

	"github.com/hypwxm/rider"
)

func RouterSub() *rider.Router {
	subrouter := rider.NewRouter()
	subrouter.GET("/midd", func(c rider.Context) {
		fmt.Println("mide in mid in mid2")
	})

	return subrouter
}
