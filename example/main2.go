package main

import (
	"rider"
	user "rider/example/router"
	"net/http"
	"fmt"
)

func main() {
	//new一个rider，创建一个app
	app := rider.New()
	app.AddMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("middleware_app")
		},
	)
	app.USE("/login", user.Router)
	app.USE("/my", user.Router)
	app.GetServer().ServerMux.HandleFunc("/se", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("server")
	})
	app.Listen(":8000")
}
