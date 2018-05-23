package main

import (
	"html/template"
	"os"
	"path/filepath"
	"rider"
)

func main() {
	app := rider.New()
	//app.ViewEngine(rider.BaseRender()) //实现BaseRender
	wd, _ := os.Getwd()
	app.SetViews(filepath.Join(wd, "src/github.com/hypwxm/rider/example/render/views"), "tpl", template.FuncMap{})
	//默认不缓存模版
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
