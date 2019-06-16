package main

import (
	"fmt"
	"time"

	"github.com/hypwxm/rider"
)

func main() {
	app := rider.New()
	app.Logger(8)
	app.USE(rider.RiderJwt("rider", "adadad", 10, &rider.RiderCookie{
		// Secure: true,
	}))
	app.GET("/token", func(c rider.Context) {
		//token, _ := c.GetLocals("token").(*rider.RiderJwter).Set("test", "test2")
		fmt.Println(c.Jwt().GetToken())
		token := c.Jwt().Get("test")
		//jwt.SetTokenCookie(c)
		//c.SetHeader("token", token)
		c.Send(200, []byte(token.(string)))
	})

	app.GET("/settoken", func(c rider.Context) {
		fmt.Println(time.Now().Unix(), time.Now().UnixNano())
		c.Jwt().Set("test", "test")
		c.Send(200, []byte(c.Jwt().GetToken()))
	})
	app.GET("/tokenparse", func(c rider.Context) {
		c.Logger().INFO(c.CookieValue("token"))
		//token := c.Jwt.Jwt.TokenString
		c.Jwt().Delete("test")
		c.Jwt().DeleteAll()
		c.Jwt().Set("a", "b")
		fmt.Println(c.Jwt().Values())
		fmt.Println(c.Jwt().Get("a"))
		fmt.Println(c.Jwt().GetToken())
	})
	app.Listen(":5002")
}
