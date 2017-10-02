package router

import (
	"fmt"
	"rider"
)


func Router() *rider.Router {
	router := rider.NewRouter()

	router.AddMiddleware(
		func(context *rider.Context) {
			fmt.Println("middleware_root")
			context.Next()
		},
	)

	router.GET("/sub", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Send([]byte("xxxx"))
		},
	})
	router.POST("/sub", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Send([]byte("post"))
		},
	})
	router.ANY("/next", Router3())

	router.AddMiddleware(
		func(context *rider.Context) {
			fmt.Println("middleware_rootp")
			context.Next()
		},
	)

	return router
}
