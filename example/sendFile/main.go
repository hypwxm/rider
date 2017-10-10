package main

import (
	"rider"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	app := rider.New()
	wd, _ := os.Getwd()
	app.GET("/file", &rider.Router{
		Handler: func(c *rider.Context) {
			files, err := ioutil.ReadDir(filepath.Join(wd, "src/rider/example/sendFile"))
			if err != nil {
				panic(err)
			}
			file := files[0]
			c.SendFile(filepath.Join(wd, "src/rider/example/sendFile", file.Name()))
		},
	})
	app.Listen(":5007")
}
