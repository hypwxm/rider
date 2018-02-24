package main

import (
	"fmt"
	"rider"
	"time"
)

func main() {
	app := rider2.New()
	app.Logger(8)
	app.USE(rider2.RiderJwt("rider", time.Hour))
	app.GET("/token", func(c rider2.Context) {
		//token, _ := c.GetLocals("token").(*rider.RiderJwter).Set("test", "test2")
		token, _ := c.Jwt().Set("test", " test")
		//jwt.SetTokenCookie(c)
		//c.SetHeader("token", token)
		c.Send(200, []byte(token))
	})
	app.GET("/tokenparse", func(c rider2.Context) {
		c.Logger().INFO(c.CookieValue("token"))
		//token := c.Jwt.Jwt.TokenString
		c.Jwt().Delete("test")
		c.Jwt().DeleteAll()
		c.Jwt().Set("a", "b")
		fmt.Println(c.Jwt().Claims())
		fmt.Println(c.Jwt().ClaimsValue("a"))
	})
	app.Listen(":5002")
}
