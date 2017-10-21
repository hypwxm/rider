package main

import (
	"rider"
	"net/http"
)

func main() {
	app := rider.New()
	app.GET("/setCookie", &rider.Router{
		Handler: func(c *rider.Context) {
			//设置cookie
			c.SetCookie(http.Cookie{
				Name: "cookie",
				Value: "rider",
				MaxAge: 100000, //秒
				Path: "/",  //默认为 "/"
			})

			c.Send([]byte("ok"))
		},
	})
	app.GET("/getCookie", &rider.Router{
		Handler: func(c *rider.Context) {
			cookie, _ := c.CookieValue("cookie")
			//cookies := c.Cookies()
			c.Send([]byte(cookie))
		},
	})
	app.GET("/removeCookie", &rider.Router{
		Handler: func(c *rider.Context) {
			//删除cookie
			c.RemoveCookie("cookie")
			c.Send([]byte("ok"))
		},
	})
	app.Listen(":5003")
}
