package main

import (
	"rider"
	"time"
)

func main() {
	app := rider.New()
	app.Logger(8)
	app.GET("/", func(c rider.Context) {
		go func() {
			<-c.Response().CloseNotify()
			c.Logger().WARNING("http closed")
		}()
		time.Sleep(20e9)
	})
	app.Listen(":5000")
}
