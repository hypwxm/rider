package main

import (
	"rider"
	"net/http"
)

func main() {
	app := rider3.New()
	app.GET("/setCookie", &rider3.Router{
		Handler: func(c *rider3.Context) {
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
	app.GET("/getCookie", &rider3.Router{
		Handler: func(c *rider3.Context) {
			cookie, _ := c.CookieValue("cookie")
			//cookies := c.Cookies()
			c.Send([]byte(cookie))
		},
	})
	app.GET("/removeCookie", &rider3.Router{
		Handler: func(c *rider3.Context) {
			//删除cookie
			c.RemoveCookie("cookie")
			c.Send([]byte("ok"))
		},
	})
	app.Listen(":5003")
}
