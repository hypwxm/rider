package router

import (
	"rider"
)

func Router3() *rider2.Router {
	router := rider2.NewRouter()

	router.GET("/sub", func(c rider2.Context) {
		c.Send(200, []byte("第二层子集"))
	})

	return router
}
