package main

import (
	"rider"
	"os"
	"path/filepath"
)

func main() {
	app := rider.New()
	//app.ViewEngine(rider.BaseRender()) //实现BaseRender
	wd, _ := os.Getwd()
	app.SetViews(filepath.Join(wd, "src/rider/example/render/views"), "tpl")
	//默认不缓存模版
	app.CacheViews()
	app.GET("/", func(c rider.Context) {
		c.Render("render", map[string]interface{}{
			"a": "b",
		})
	})
	app.GET("/w", func(c rider.Context) {
		c.Render("worker/w", map[string]interface{}{
			"json": map[string]int{
				"a": 1,
				"b": 2,
			},
		})
	})
	app.Listen(":5010")
}