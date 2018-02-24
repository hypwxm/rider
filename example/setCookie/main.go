package main

import (
	"net/http"
	"rider"
)

func main() {
	app := rider2.New()
	app.GET("/setCookie", func(c rider2.Context) {
		//设置cookie
		c.SetCookie(http.Cookie{
			Name:   "cookie",
			Value:  "rider",
			MaxAge: 100000, //秒
			Path:   "/",    //默认为 "/"
		})

		c.Send(200, []byte("ok"))
	})
	app.GET("/getCookie", func(c rider2.Context) {
		cookie, _ := c.CookieValue("cookie")
		//cookies := c.Cookies()
		c.Send(200, []byte(cookie))
	})
	app.GET("/removeCookie", func(c rider2.Context) {
		//删除cookie
		c.DeleteCookie("cookie")
		c.Send(200, []byte("ok"))
	})
	app.Listen(":5003")
}
