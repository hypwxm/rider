package router

import (
	"fmt"
	"rider"
)


func Router3() *rider.Router {
	router := rider.NewRouter()

	router.AddMiddleware(
		func(context *rider.Context) {
			fmt.Println("middleware_root")
			context.Next()
		},
	)

	router.GET("/sub", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Response.Send("xxxx")
		},
	})
	router.AddMiddleware(
		func(context *rider.Context) {
			fmt.Println("middleware_root2")
			context.Next()
		},
	)
	router.POST("/sub", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Response.Send("post")
		},
	})



	return router
}
