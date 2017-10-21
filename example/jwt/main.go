package main

import (
	"rider"
	"time"
	"fmt"
)

func main() {
	app := rider.New()
	app.Logger(8)
	app.AddMiddleware(rider.RiderJwt("rider", time.Hour))
	app.GET("/token", &rider.Router{
		Handler: func(c *rider.Context) {
			//token, _ := c.GetLocals("token").(*rider.RiderJwter).Set("test", "test2")
			token, _ := c.Jwt.Set("test", " test")
			//jwt.SetTokenCookie(c)
			//c.SetHeader("token", token)
			c.Send([]byte(token))
		},
	})
	app.GET("/tokenparse", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Logger.INFO(c.CookieValue("token"))
			//token := c.Jwt.Jwt.TokenString
			c.Jwt.Delete("test")
			c.Jwt.DeleteAll()
			c.Jwt.Set("a", "b")
			fmt.Println(c.Jwt.Claims())
			fmt.Println(c.Jwt.ClaimsValue("a"))
			c.End()
		},
	})
	app.Listen(":5002")
}