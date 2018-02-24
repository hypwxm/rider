package router

import (
	"rider"
)

func Router() *rider2.Router {
	router := rider2.NewRouter()

	router.GET("/sub", func(c rider2.Context) {
		c.Send(200, []byte("第一层子集"))
	})
	router.Kid("/next", Router3())

	return router
}
