package main

import (
	"rider"
	"time"
)

func main() {
	app := rider3.New()
	app.Logger(8)
	app.GET("/", &rider3.Router{
		Handler: func(c *rider3.Context) {
			go func() {
				<- c.CloseNotify()
				c.Logger.WARNING("http closed")
			}()
			time.Sleep(20e9)
		},
	})
	app.Listen(":5000")
}
