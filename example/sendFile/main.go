package main

import (
	"rider"
	"os"
	"path/filepath"
)

func main() {
	app := rider.New()
	app.Logger(8)
	wd, _ := os.Getwd()
	app.GET("/file", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Hijack()
			//c.Send([]byte("xx"))
			c.SendFile(filepath.Join(wd, "src/rider/example/sendFile/9f26d8773445ac485663e91080c75877.html"))
		},
	})
	app.Listen(":5007")
}
