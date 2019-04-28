package main

import (
	"fmt"
	"time"

	"github.com/hypwxm/rider"
)

func main() {
	app := rider.New()
	app.Logger(8)
	app.USE(rider.RiderJwt("rider", time.Hour))
	app.GET("/token", func(c rider.Context) {
		//token, _ := c.GetLocals("token").(*rider.RiderJwter).Set("test", "test2")
		token, _ := c.Jwt().Set("test", " test")
		//jwt.SetTokenCookie(c)
		//c.SetHeader("token", token)
		c.Send(200, []byte(token))
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
