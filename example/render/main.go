package main

import (
	"os"
	"path/filepath"
	"rider"
)

func main() {
	app := rider2.New()
	//app.ViewEngine(rider.BaseRender()) //实现BaseRender
	wd, _ := os.Getwd()
	app.SetViews(filepath.Join(wd, "src/rider/example/render/views"), "tpl")
	//默认不缓存模版
	app.CacheViews()
	app.GET("/", func(c rider2.Context) {
		c.Render("render", map[string]interface{}{
			"a": "b",
		})
	})
	app.GET("/w", func(c rider2.Context) {
		c.Render("worker/w", map[string]interface{}{
			"json": map[string]int{
				"a": 1,
				"b": 2,
			},
		})
	})
	app.Listen(":5010")
}
