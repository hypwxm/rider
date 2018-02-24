package router

import (
	"github.com/hypwxm/rider"
)

func Router() *rider.Router {
	router := rider.NewRouter()

	router.GET("/sub", func(c rider.Context) {
		c.Send(200, []byte("第一层子集"))
	})
	router.Kid("/next", Router3())

	return router
}
