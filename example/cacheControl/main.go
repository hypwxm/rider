package main

import (
	"rider"
	"rider/modules"
	"net/http"
	"strconv"
	"os"
	"path/filepath"
)

func main() {
	app := rider3.New()
	app.USE(modules.CacheControl(6000000))
	wd, _ := os.Getwd()
	var a = 0
	app.GET("/cache", &rider3.Router{
		Handler: func(context rider3.Context) {
			context.Send(http.StatusOK, []byte(strconv.Itoa(a)))
			a++
		},
	})
	app.GET("/cache2", &rider3.Router{
		Handler: func(context rider3.Context) {
			err := context.SendFile(filepath.Join(wd, "src/rider/example/cacheControl/main.go"))
			context.Logger().ERROR(err)
		},
	})
	app.Listen(":5003")
}
