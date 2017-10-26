package main

import (
	"rider"
	"os"
	"path/filepath"
	"time"
)

func main() {
	app := rider3.New()
	wd, _ := os.Getwd()
	app.SetViews(filepath.Join(wd, "src/rider/example/uploadFile/views"), "tpl")
	//app.CacheViews()
	app.GET("/up1", &rider3.Router{
		Handler: func(c *rider3.Context) {
			c.Hijack()
			c.Render("uploads", nil)
			time.Sleep(10e9)

		},
	})
	app.GET("/up2", &rider3.Router{
		Handler: func(c *rider3.Context) {
			c.Render("uploadsFiles", nil)
		},
	})
	//当文件上传
	app.POST("/uploads", &rider3.Router{
		Handler: func(c *rider3.Context) {
			formFile, err := c.FormFile("file")
			if err != nil {
				panic(err)
			}
			c.StoreFormFile(formFile, filepath.Join(wd, "src/rider/example/uploadFile", formFile.Name))
			c.SendFile(filepath.Join(wd, "src/rider/example/uploadFile", formFile.Name))
		},
	})
	//多文件上传
	app.POST("/uploadsFiles", &rider3.Router{
		Handler: func(c *rider3.Context) {
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
