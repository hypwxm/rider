package router

import (
	"rider"
)


func Router3() *rider3.Router {
	router := rider3.NewRouter()

	router.GET("/sub", &rider3.Router{
		Handler: func(c *rider3.Context) {
			c.Send([]byte("第二层子集"))
		},
	})

	return router
}
