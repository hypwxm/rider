package main

import (
	"rider"
	"os"
	"path/filepath"
	"time"
)

func main() {
	app := rider.New()
	wd, _ := os.Getwd()
	app.SetViews(filepath.Join(wd, "src/rider/example/uploadFile/views"), "tpl")
	//app.CacheViews()
	app.GET("/up1", func(c rider.Context) {
			c.Hijack()
			c.Render("uploads", nil)
			time.Sleep(10e9)
	})
	app.GET("/up2", func(c rider.Context) {
			c.Render("uploadsFiles", nil)
	})
	//当文件上传
	app.POST("/uploads", func(c rider.Context) {
			formFile, err := c.FormFile("file")
			if err != nil {
				panic(err)
			}
			c.StoreFormFile(formFile, filepath.Join(wd, "src/rider/example/uploadFile", formFile.Name))
			c.SendFile(filepath.Join(wd, "src/rider/example/uploadFile", formFile.Name))
	})
	//多文件上传
	app.POST("/uploadsFiles", func(c rider.Context) {
			formFiles, err := c.FormFiles("files")
			if err != nil {
				panic(err)
			}
			c.StoreFormFiles(formFiles, filepath.Join(wd, "src/rider/example/uploadFile"))
			c.Send(200, []byte("ok"))
	})
	app.Listen(":5002")
}
