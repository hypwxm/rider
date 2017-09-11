package router

import (
	"rider/riderRouter"
	"fmt"
	"net/http"
)


func Router() *riderRouter.Router {
	router := riderRouter.NewRouter(nil)

	/*Router.Middleware = riderRouter.MiddleWare(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("middleware_root")
		},
	)*/


	router.AddMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("middleware_root")
		},
	)

	/*route := riderRouter.NewRouter(nil)
	route.AddMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("middleware1")
		},
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("middleware2")
		},
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("middleware3")
		},
	)

	Router.USE("/middle", route)*/

	//var User = router.NewRouter(nil)
	router.POST("/user", &riderRouter.Router{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("xxx")
		},
		Middleware: riderRouter.MiddleWare(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("middleware1")
			},
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("middleware2")
			},
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("middleware3")
			},
		),
	})
	router.ANY("/kong", &riderRouter.Router{})

	router.AddMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("next")
		},
	)


	router.GET("/getuser", &riderRouter.Router{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("zzz")
		},
	})


	return router
}
