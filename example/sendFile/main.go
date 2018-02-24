package main

import (
	"os"
	"path/filepath"
	"rider"
)

func main() {
	app := rider2.New()
	app.Logger(8)
	wd, _ := os.Getwd()
	app.GET("/file", func(c rider2.Context) {
		c.Hijack()
		//c.Send([]byte("xx"))
		c.SendFile(filepath.Join(wd, "src/rider/example/sendFile/9f26d8773445ac485663e91080c75877.html"))
	})
	app.Listen(":5007")
}
