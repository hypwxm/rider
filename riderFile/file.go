package riderFile

import (
	"os"
	"path/filepath"
	"log"
	"strings"
)

//get filename extensions
func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && path[i] != '/'; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}
	return ""
}


//判断文件或者目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

//判断是否为目录
func IsDir(path string) bool {
	if IsExist(path) {
		fi, _ := os.Stat(path)
		return fi.IsDir()
	}
	return false
}

//获取当前目录
func GetCWD() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}