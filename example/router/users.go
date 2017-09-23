package router

import (
	"fmt"
	"rider"
)


func Router() *rider.Router {
	router := rider.NewRouter(nil)

	/*Router.Middleware = riderRouter.MiddleWare(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("middleware_root")
		},
	)*/


	router.AddMiddleware(
		func(context *rider.Context) {
			fmt.Println("middleware_root")
			context.Next()
		},
	)

	route := rider.NewRouter(nil)


	router.ANY("/middle", Router2())
	route.GET("/sub", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Response.Send("xxxx")
		},
	})

	//var User = router.NewRouter(nil)
	router.POST("/user", &rider.Router{
		Handler: func(context *rider.Context) {
			fmt.Println("xxx")
		},
		Middleware: rider.MiddleWare(
			func(context *rider.Context) {
				fmt.Println("middleware1")
				context.Next()
			},
			func(context *rider.Context) {
				fmt.Println("middleware2")
				context.Next()
			},
			func(context *rider.Context) {
				context.Next()
				return
				fmt.Println("middleware3")
			},
		),
	})
	router.ANY("/kong", &rider.Router{})

	router.AddMiddleware(
		func(context *rider.Context) {
			fmt.Println("next")
		},
	)


	router.GET("/getuser", &rider.Router{
		Handler: func(context *rider.Context) {
			fmt.Println("zzz")
		},
	})


	return router
}
