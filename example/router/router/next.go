package router

import (
	"rider"
)


func Router3() *rider.Router {
	router := rider.NewRouter()

	router.GET("/sub", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Send([]byte("第二层子集"))
		},
	})

	return router
}
