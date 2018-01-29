package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"rider/utils/file"
	"strconv"
	"strings"
	"time"
)

//为日志存储设置多功能存储（默认终端，输出到日志文件，发送到邮箱）
//日志文件分两种（一种记录fatal，panic，error，warning等级的日志；一种记录info，console等级的日志）
//日志文件的创建以当前时间点为前缀，参数加入服务进程
//日志文件默认最大为20M，大于20M后会创建新的日志文件。旧的日志文件会被压缩（设置成可配置）

func (lq *LogQueue) SetLogOut(out io.Writer) {
	lq.logWriter = log.New(out, "", 0)
}
func (lq *LogQueue) SetErrLogOut(out io.Writer) {
	lq.errLogWriter = log.New(out, "", 0)
}

//生成日志文件
func createLogFile(filename string) (*os.File, error) {
	now := time.Now().Format("2006-01-02T15:04:05")
	pidStr := strconv.Itoa(os.Getpid())
	logFileName := fmt.Sprintf("%s@%s.txt", now, pidStr)
	lf, err := os.OpenFile(filepath.Join(filename, logFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return lf, nil
}

//生成错误日志文件
func createErrLogFile(filename string) (*os.File, error) {
	now := time.Now().Format("2006-01-02T15:04:05")
	pidStr := strconv.Itoa(os.Getpid())
	errlogFileName := fmt.Sprintf("%s@%s_error.txt", now, pidStr)
	elf, err := os.OpenFile(filepath.Join(filename, errlogFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return elf, nil
}

//设置输出到日志文件
//默认输出到服务的当前目录
//如果filename已存在并且为文件夹，则会往里添加文件
//如果filename已存在但是是一个文件，则会在同级目录下生成日志文件
//新建一个名为filename的文件夹
//日志文件创建失败，会将日志转移至os.Stdout
func (lq *LogQueue) SetLogOutPath(filename string) error {
	if strings.TrimSpace(filename) == "" {
		filename = file.GetCWD()
	}
	//判断filename是否存在，并且不为文件夹
	if file.IsExist(filename) && !file.IsDir(filename) {
		filename = filepath.Dir(filepath.Clean(filename))
	} else {
		err := os.MkdirAll(filepath.Dir(filename), 0777)
		if err != nil {
			return err
		}
	}
	lq.logOutPath = filename
	lf, err := createLogFile(filename)
	if err != nil {
		lq.RemoveDestination(1)
		lq.PANIC(err)
		return err
	}
	elf, err := createErrLogFile(filename)
	if err != nil {
		lq.RemoveDestination(1)
		lf.Close()
		lq.PANIC(err)
		return err
	}
	lq.SetLogOut(lf)
	lq.SetErrLogOut(elf)
	lq.logOutFile = lf.Name()
	lq.errLogOutFile = elf.Name()
	lq.logOutFile_FD = lf
	lq.errLogOutFile_FD = elf
	lq.AddDestination(1)
	return nil
}

//打开日志文件，获取日志文件描述符，（一般用户日志文件描述符以外关闭）
func (lq *LogQueue) getFileFd(fileName string) (*os.File, error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//设置日志文件大小上线
func (lq *LogQueue) SetLogFileMaxSize(size int64) {
	lq.maxLogFileSize = size
}

//判断日志文件大小，并判断是否需要生成新的日志文件(没生成文件，或者filename不存在，则返回旧文件，获取返回的新文件的filename)
func (lq *LogQueue) updateLogFile(filename string) error {
	f, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return errors.New(filename + " is a directory")
	}
	filename = filepath.Dir(filename)
	if f.Size() >= lq.maxLogFileSize {
		lf, err := createLogFile(filename)
		if err != nil {
			return err
		}
		lq.SetLogOut(lf)
		lq.logOutFile_FD = lf
		lq.logOutFile = lf.Name()
	}
	return nil
}

//更新错误日志文件
func (lq *LogQueue) updateErrLogFile(filename string) error {
	f, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return errors.New(filename + " is a directory")
	}
	filename = filepath.Dir(filename)
	if f.Size() >= lq.maxLogFileSize {
		lf, err := createErrLogFile(filename)
		if err != nil {
			return err
		}
		lq.SetErrLogOut(lf)
		lq.errLogOutFile_FD = lf
		lq.errLogOutFile = lf.Name()
	}
	return nil
}
