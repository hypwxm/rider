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
	app.AddMiddleware(rider.Gzip(-1))
	app.SetStatic(filepath.Join(wd, "src/rider/example/setStatic/public"))
	app.GET("/xxx", &rider.Router{
		Handler: func(c *rider.Context) {

			//time.Sleep(15e9)
			//c.Redirect(307, "https://www.baidu.com")
			//panic(errors.New("xxx"))
			c.Hijack()
			c.Send(200, []byte("点滴家啊阿娇无敌啊叫i 啊叫的i 啊基地是奥会降低啊阿娇的i 哦啊叫"))
		},
	})
	app.Listen(":5001")
}
