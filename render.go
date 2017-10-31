package rider

import (
	"io/ioutil"
	"os"
	"html/template"
	"fmt"
	"github.com/hypwxm/rider/riderFile"
	"io"
	"errors"
	"path/filepath"
)

var (
	//tplsRender BaseRender = newRender() //默认自带模板引擎
)

type BaseRender interface {
	Render(w io.Writer, tplName string, data interface{}) error  //tplName模板名称 ,data模板数据
}

var _ BaseRender = newRender()

type render struct {
	templates map[string]*template.Template
	//server *HttpServer
	tplDir string
	extName string
	cache bool
}

func newRender() *render {
	return &render{
		templates: make(map[string]*template.Template),
		cache: false,
	}
}

var (
	templateName string
	templatePath string
	fullTmplPath string
	nextPrefix string
)

func (rd *render) registerTpl(tplDir string, extName string, namePrefix string) {
	fileInfoArr, err := ioutil.ReadDir(tplDir)

	if err != nil {
		panic(err)
	}

	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		templatePath = tplDir + "/" + templateName
		f, err := os.Stat(templatePath)
		if err != nil {
			panic(err)
		}
		isdir := f.IsDir()
		if isdir {
			if namePrefix == "" {
				nextPrefix = templateName
			} else {
				nextPrefix = namePrefix + "/" + templateName
			}

			rd.registerTpl(tplDir + "/" + templateName, extName, nextPrefix)
			continue
		}

		if ext := riderFile.Ext(templateName); ext != extName {
			continue
		}

		templateName = templateName[:(len(templateName) - len(extName) - 1)]
		//templatePaths = append(templatePaths, templatePath)
		t := template.Must(template.ParseFiles(templatePath))

		if namePrefix == "" {
			fullTmplPath = templateName
		} else {
			fullTmplPath = namePrefix + "/" + templateName
		}

		rd.templates[fullTmplPath] = t
	}
	if GlobalENV == ENV_Development || GlobalENV == ENV_Debug {
		fmt.Println("templates was loaded over!")
	}
}

func (rd *render) Render(w io.Writer, tplName string, data interface{}) error {
	//如果设置模板缓存，则从缓存中读取模板
	if rd.cache {
		if views, ok := rd.templates[tplName]; ok {
			views.Execute(w, data)
			return nil
		}
		return errors.New("未找到" + tplName + "模板信息")
	} else {
		//如果没设置缓存，需要从disk中直接读取模板文件
		tplPath := filepath.Join(rd.tplDir, tplName + "." + rd.extName)
		f, err := os.Stat(tplPath)
		if err != nil {
			return err
		}
		if f.IsDir() {
			return errors.New(tplName + "模板是个目录")
		}
		t := template.Must(template.ParseFiles(tplPath))
		return t.Execute(w, data)
	}
}

func (rd *render) Cache() {
	rd.cache = true
}

func (rd *render) setTplDir (tplDir string) {
	rd.tplDir = tplDir
}

func (rd *render) setExtName (extName string) {
	rd.extName = extName
}