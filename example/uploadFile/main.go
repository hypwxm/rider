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
	app.GET("/up1", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Hijack()
			c.Render("uploads", nil)
			time.Sleep(10e9)

		},
	})
	app.GET("/up2", &rider.Router{
		Handler: func(c *rider.Context) {
			c.Render("uploadsFiles", nil)
		},
	})
	//当文件上传
	app.POST("/uploads", &rider.Router{
		Handler: func(c *rider.Context) {
			formFile, err := c.FormFile("file")
			if err != nil {
				panic(err)
			}
			c.StoreFormFile(formFile, filepath.Join(wd, "src/rider/example/uploadFile", formFile.FileName()))
			c.SendFile(filepath.Join(wd, "src/rider/example/uploadFile", formFile.FileName()))
		},
	})
	//多文件上传
	app.POST("/uploadsFiles", &rider.Router{
		Handler: func(c *rider.Context) {
			formFiles, err := c.FormFiles("files")
			if err != nil {
				panic(err)
			}
			c.StoreFormFiles(formFiles, filepath.Join(wd, "src/rider/example/uploadFile"))
			c.Send([]byte("ok"))
		},
	})
	app.Listen(":5002")
}