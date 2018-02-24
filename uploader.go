package rider

//思路来源，大神dotweb的框架：更多框架信息=>

import (
	"errors"
	"github.com/hypwxm/rider/utils/file"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type UploadFile struct {
	File   multipart.File
	header *multipart.FileHeader
	Ext    string
	Name   string
	size   int64
}

func NewUploadFile(f multipart.File, header *multipart.FileHeader) *UploadFile {
	return &UploadFile{
		File:   f,
		header: header,
		Name:   header.Filename,
		Ext:    file.Ext(header.Filename),
	}
}

// 获取文件大小的接口
type size interface {
	Size() int64
}

//获取上传的文件大小，具体需要的时候在进行赋值
func (f *UploadFile) Size() int64 {
	if f.size <= 0 {
		if sizer, ok := f.File.(size); ok {
			f.size = sizer.Size()
		}
	}
	return f.size
}

//存储http请求的文件到服务器
func (f *UploadFile) StoreFile(filename string) (size int64, err error) {
	if strings.TrimSpace(filename) == "" {
		return 0, errors.New("filename cannot be empty")
	}
	fw, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return io.Copy(fw, f.File)
}
