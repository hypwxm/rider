package router

import (
	"rider"
)


func Router() *rider.Router {
	router := rider.NewRouter()

	router.GET("/sub", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Send([]byte("第一层子集"))
		},
	})
	router.GET("/next", Router3())


	return router
}
