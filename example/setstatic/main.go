package main

import (
	"rider"
	"os"
	"path/filepath"
)

func main() {
	app := rider3.New()
	app.Logger(8)

	wd, _ := os.Getwd()
	app.USE(rider3.Gzip(-1))
	app.SetStatic(filepath.Join(wd, "src/rider/example/setStatic/public"))
	app.GET("/xxx", &rider3.Router{
		Handler: func(c rider3.Context) {

			//time.Sleep(15e9)
			//c.Redirect(307, "https://www.baidu.com")
			//panic(errors.New("xxx"))
			c.Hijack()
			c.Logger().DEBUG("xxx")
			c.Send(200, []byte("点滴家啊阿娇无敌啊叫i 啊叫的i 啊基地是奥会降低啊阿娇的i 哦啊叫"))
		},
	})
	app.Listen(":5001")
}
