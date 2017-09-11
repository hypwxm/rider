package router2

import (
	"rider/riderRouter"
	"fmt"
	"net/http"
)

var Router = riderRouter.NewRouter(nil)

func init() {

	/*Router.Middleware = riderRouter.MiddleWare(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("middleware_root")
		},
	)*/

	Router.AddMiddleware(
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
	Router.POST("/user", &riderRouter.Router{
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
	Router.USE("/kong", &riderRouter.Router{})

	Router.AddMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("next")
		},
	)


	Router.GET("/getuser", &riderRouter.Router{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("zzz")
		},
	})



}
