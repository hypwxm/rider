package main

import (
	"rider"
	"os"
	"path/filepath"
	"time"
	"fmt"
)

func main() {
	app := rider.New()
	app.Logger(8)
	wd, _ := os.Getwd()
	app.AddMiddleware(rider.Gzip(-1))
	app.SetStatic(filepath.Join(wd, "src/rider/example/setStatic/public"))
	app.GET("/xx", &rider.Router{
		Handler: func(c *rider.Context) {
			start := time.Now()
			c.Hijack()
			fmt.Println(time.Now().Sub(start))

			c.Send([]byte("adad"))
		},
	})
	app.Listen(":5001")
}
