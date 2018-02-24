package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"rider"
)

func main() {
	app := rider2.New()
	app.Logger(8)
	wd, _ := os.Getwd()
	app.GET("/download", func(c rider2.Context) {
		files, _ := ioutil.ReadDir(filepath.Join(wd, "src/rider/example/download"))
		file := files[0]
		c.Download(filepath.Join(wd, "src/rider/example/download", file.Name()), "a   dad", "attachment")
	})
	app.Listen(":5004")
}
