package main

import (
	"rider"
)

func main() {
	app := rider.New()
	rlog := app.Logger(8)
	app.GET("/logger", &rider.Router{
		Handler: func(c *rider.Context) {
			rlog.INFO("xx", "yy")
			rlog.DEBUG("OK")
			rlog.PANIC("adad")
			c.SendJson(map[string]string{
				"a":"1",
			})
		},
	})
	app.GET("/log500", &rider.Router{
		Handler: func(c *rider.Context) {
			c.SetStatusCode(500)
			c.Send([]byte("error500"))
		},
	})
	app.GET("/300", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Hijack()
			c.SetStatusCode(300)
			//panic(errors.New("adad"))
			c.Send([]byte("300"))
			c.Send([]byte("awd"))
			//c.SetStatusCode(400)
		},
	})
	app.GET("/ada", &rider.Router{})
	app.Listen(":5000")

}



