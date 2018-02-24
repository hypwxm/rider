package main

import (
	"net/http"
	"os"
	"path/filepath"
	"rider"
	"rider/modules"
	"strconv"
)

func main() {
	app := rider2.New()
	app.USE(modules.CacheControl(6000000))
	wd, _ := os.Getwd()
	var a = 0
	app.GET("/cache", func(context rider2.Context) {
		context.Send(http.StatusOK, []byte(strconv.Itoa(a)))
		a++
	})
	app.GET("/cache2", func(context rider2.Context) {
		err := context.SendFile(filepath.Join(wd, "src/rider/example/cacheControl/main.go"))
		context.Logger().ERROR(err)
	})
	app.Listen(":5003")
}
