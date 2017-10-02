package main

import (
	"rider"
	user "rider/example/router"
	//user2 "rider/example/router2"
	"fmt"
	"flag"
	"net/http"
)

type pint int

func crocess(c *rider.Context) {
	fmt.Println("middleware_app")

	c.SetHeader("Access-Control-Allow-Origin", "*")
	//c.Send([]byte("xxx"))
	//return
	//time.Sleep(10e9)
	c.SetLocals("locals", "adada")
	c.Next()
}

func main() {

	env := flag.String("env", rider.ENV_Development, "设置app环境变量")
	flag.Parse()

	rider.SetEnvMode(*env)

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


	/*app.POST("/haa", &rider.Router{
		Handler:func(c *rider.Context) {
			//c.Response.Send(`{"a":"1"}`)
			//c.Response.SetContentType("application/json")

			panic(errors.New("adadad"))

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
	})*/

	app.GET("/xx/:id/:id2", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Hijack()

			c.CookieValue("xxx")

			c.ResCookie(http.Cookie{
				Name: "xxx",
				Value: "yyy=awdad",
				Path: "/",  // optional
				MaxAge: 100,
			})
			c.ResCookie(http.Cookie{
				Name: "ad",
				Value: "yyy",
				Path: "/",  // optional
				MaxAge: -1,
			})
			fmt.Println(c.ResponseHeader())

			fmt.Println("xxx")
			c.Send([]byte("xx"))
			//fmt.Println(file[0])
		},
	})

	app.ANY("/login", user.Router())

	//app.ANY("/ha", user.Router())

	/*app.AddMiddleware(
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
	})*/
	app.Listen(":8000")
}
