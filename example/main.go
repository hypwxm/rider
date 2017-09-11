package main

import (
	"rider"
	user "rider/example/router"
	//user2 "rider/example/router2"

	"net/http"
	"fmt"
	"rider/riderRouter"
)

func main() {
	//new一个rider，创建一个app
	app := rider.New()
	app.AddMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("middleware_app")
		},
	)
	/*app.GET("/haa", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("haaaa")
	})*/
	app.ANY("/haa", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("xxxx")
	})

	app.GET("/haa", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("yyy")
	})

	app.ANY("/login", user.Router())
	app.AddMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("appnext")
		},
	)

	app.ANY("/my", user.Router())

	app.HandleFunc("/xx", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("server")
	})


	app.GET("/justrouter", &riderRouter.Router{
		Handler: func(w http.ResponseWriter, R *http.Request) {
			fmt.Println("justrouter")
		},
		Middleware: riderRouter.MiddleWare(
			func(w http.ResponseWriter, R *http.Request) {
				fmt.Println("justrouter_middleware")
			},
		),
	})
	app.Listen(":8000")
}
