package router

import (
	"rider"
)


func Router3() *rider.Router {
	router := rider.NewRouter()

	router.GET("/sub", func(c rider.Context) {
			c.Send(200, []byte("第二层子集"))
	})

	return router
}
