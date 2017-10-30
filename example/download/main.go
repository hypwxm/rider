package main

import (
	"rider"
	"io/ioutil"
	"path/filepath"
	"os"
)

func main() {
	app := rider.New()
	app.Logger(8)
	wd, _ := os.Getwd()
	app.GET("/download", func(c rider.Context) {
			files, _ := ioutil.ReadDir(filepath.Join(wd, "src/rider/example/download"))
			file := files[0]
			c.Download(filepath.Join(wd, "src/rider/example/download", file.Name()), "a   dad", "attachment")
	})
	app.Listen(":5004")
}
