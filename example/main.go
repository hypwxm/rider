package main

import (
	"rider"
	user "rider/example/router"
	//user2 "rider/example/router2"
	"fmt"
)

type pint int

func crocess(c *rider.Context) {
	fmt.Println("middleware_app")
	c.Response.SetHeader("Access-Control-Allow-Origin", "*")
	c.Next()
}

func main() {



	//new一个rider，创建一个app
	app := rider.New()
	app.AddMiddleware(
		crocess,
	)
	/*app.GET("/haa", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("haaaa")
	})*/


	type st struct {
		A string
		b int
	}


	app.POST("/haa", &rider.Router{
		Handler:func(c *rider.Context) {
			//c.Response.Send(`{"a":"1"}`)
			//c.Response.SetContentType("application/json")

			x := c.Request.QueryStrings()
			fmt.Println(x["xx"] == nil)
			//wd, _ := os.Getwd()
			//var a = 131313
			c.Response.SetStatusCode(201)

			//f, _ := os.Open(filepath.Join(wd, "src", "rider", "serve.go"))
			c.Response.Send(map[string]interface{}{
				"a": map[string]interface{}{
					"a":[]int{1,2,3},
				},
			})

			//c.Response.Send(a)

		},
	})

	app.ANY("/login", user.Router())
	app.AddMiddleware(
		func(context *rider.Context) {
			fmt.Println("appnext")
			context.Next()
		},
	)

	app.ANY("/my", user.Router())

	app.GET("/justrouter", &rider.Router{
		Handler: func(context *rider.Context) {
			fmt.Println("justrouter")
		},
		Middleware: rider.MiddleWare(
			func(context *rider.Context) {
				fmt.Println("justrouter_middleware")
			},
		),
	})
	app.Listen(":8000")
}
