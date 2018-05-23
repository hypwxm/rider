package rider

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/hypwxm/rider/utils/file"
)

var (
//tplsRender BaseRender = newRender() //默认自带模板引擎
)

type BaseRender interface {
	Render(w io.Writer, tplName string, data interface{}) error //tplName模板名称 ,data模板数据
}

var _ BaseRender = newRender()

type render struct {
	templates *template.Template
	tplDir    string
	extName   string
	FuncMap   template.FuncMap
}

func newRender() *render {
	return &render{
		templates: template.New("app").Delims("{%", "%}"),
		FuncMap:   template.FuncMap{},
	}
}

var (
	templateName string
	templatePath string
	fullTmplPath string
	nextPrefix   string
)

func (rd *render) registerTpl(tplDir string, extName string, funcMap template.FuncMap, namePrefix string) *render {
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

			rd.registerTpl(tplDir+"/"+templateName, extName, funcMap, nextPrefix)
			continue
		}

		if ext := file.Ext(templateName); ext != extName {
			continue
		}

		templateName = templateName[:(len(templateName) - len(extName) - 1)]
		if namePrefix == "" {
			fullTmplPath = templateName
		} else {
			fullTmplPath = namePrefix + "/" + templateName
		}

		tplByte, err := ioutil.ReadFile(templatePath)
		if err != nil {
			panic(err)
		}

		parseByte := defineTemp(fullTmplPath, tplByte)
		_, err = rd.templates.Funcs(funcMap).Parse(string(parseByte))
		if err != nil {
			fmt.Println(string(parseByte))
			panic(err)
		}
	}
	if GlobalENV == ENV_Development || GlobalENV == ENV_Debug {
		fmt.Println("templates was loaded over!")
	}
	return rd
}

//实现BaseRender的Render
func (rd *render) Render(w io.Writer, tplName string, data interface{}) error {
	//如果设置模板缓存，则从缓存中读取模板
	if views := rd.templates.Lookup(tplName); views != nil {
		err := views.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
		return nil
	}
	return errors.New("未找到" + tplName + "模板信息")
}

//定义模板，返回的数据是  `{{define "tplName"}}html模板字符串{{end}}`  的字节数组
func defineTemp(tplName string, tplByte []byte) []byte {
	prefixByte := []byte(`{%define "` + tplName + `"%}`)
	suffixByte := []byte(`{%end%}`)
	preSufLen := len(append(prefixByte, suffixByte...))
	tplLen := len(tplByte) + preSufLen
	parseByte := make([]byte, tplLen)
	parseByte = append(prefixByte, append(tplByte, suffixByte...)...)
	return parseByte
}
