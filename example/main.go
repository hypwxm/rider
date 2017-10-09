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

	//time.Sleep(10e9)
	c.SetHeader("Access-Control-Allow-Origin", "*")
	//c.Send([]byte("xxx"))
	//return
	//time.Sleep(10e9)
	//panic(errors.New("adwadad"))
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

	app.GET("/", &rider.Router{

	})

	app.GET("/xx/:id/:id2", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Hijack()

			//c.Redirect(306, "http://mv.51mzzk.com")
			c.CookieValue("xxx")
			fmt.Println(c.SendHeader())

			c.SendCookie(http.Cookie{
				Name:   "xxx",
				Value:  "yyy=123",
				//Path:   "/", // optional
				MaxAge: 100,
			})

			c.AddHeader("ada", "ada")

			c.SendCookie(http.Cookie{
				Name:   "ad",
				Value:  "100",
				//Path:   "/", // optional
				MaxAge: 100,
				Domain: "mv.51mzzk.com",
			})
			c.RemoveCookie("ad")
			fmt.Println(c.SendHeader())

			//fmt.Println(c.ResponseHeader())

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
		Middleware: rider.MiddleWare(			c.Hijack()

			func(context *rider.Context) {
				fmt.Println("justrouter_middleware")
			},
		),
	})*/

	app.Error(func(c *rider.Context, err string, code int) {
		fmt.Println("ccccc")
		c.Send([]byte(err))
	})

	app.Listen(":8000")
}
