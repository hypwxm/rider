package rider

import (
	"html/template"
	"os"
	"path/filepath"
	"testing"
)

var app = New()

func TestRegisterTpl(t *testing.T) {
	wd, _ := os.Getwd()
	if tplsRender, ok := app.GetServer().tplsRender.(*render); ok {
		tplsRender.registerTpl(filepath.Join(wd, "test/views"), "tpl", template.FuncMap{}, "")
		/*if len(tplsRender.templates) != 2 {
			t.Errorf("%s", "注册的模板数量与实际目录中的符合条件的文件数量不同")
		}
		for k, _ := range tplsRender.templates {
			if k != "render" && k != "in/in" {
				t.Errorf("%v, %v", k, "逻辑错误，模板文件变量名和逻辑思路不同")
			}
		}*/
	} else {
		t.Errorf("%s", "render未实现BaseRender接口")
	}
}

func TestRender(t *testing.T) {
	wd, _ := os.Getwd()
	app.SetViews(filepath.Join(wd, "test/views"), "tpl", template.FuncMap{})

	if tplsRender, ok := app.GetServer().tplsRender.(*render); ok {
		//未做模板缓存的情况下
		//tplsRender.SetViews(filepath.Join(wd, "test/views"), "tpl", "")
		err := tplsRender.Render(os.Stdout, "render", nil)
		if err != nil {
			t.Errorf("%s, 模板路径为%s", err, filepath.Join(wd, "test/views"))
		}
		err = tplsRender.Render(os.Stdout, "in/in", nil)
		if err != nil {
			t.Errorf("%s, 模板路径为%s", err, filepath.Join(wd, "test/views"))
		}

		//模板缓存之后
		// app.CacheViews()
		err = tplsRender.Render(os.Stdout, "render", nil)
		if err != nil {
			t.Errorf("%s, 模板路径为%s", err, filepath.Join(wd, "test/views"))
		}
		err = tplsRender.Render(os.Stdout, "in/in", nil)
		if err != nil {
			t.Errorf("%s, 模板路径为%s", err, filepath.Join(wd, "test/views"))
		}
		err = tplsRender.Render(os.Stdout, "in/inn", nil)
		if err == nil {
			t.Errorf("%s, 模板路径为%s", err, filepath.Join(wd, "test/views"))
		}
		/*if len(tplsRender.templates) != 2 {
			t.Errorf("模板数量", len(tplsRender.templates), "!=", 2)
		}*/

	} else {
		t.Errorf("%s, %T", "render未实现BaseRender接口", app.GetServer().tplsRender)
	}
}
