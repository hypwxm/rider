package router

import (
	"rider"
)


func Router() *rider3.Router {
	router := rider3.NewRouter()

	router.GET("/sub", &rider3.Router{
		Handler: func(c *rider3.Context) {
			c.Send([]byte("第一层子集"))
		},
	})
	router.GET("/next", Router3())


	return router
}
