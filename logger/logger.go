package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rider/smtp/FlyWhisper"
	"rider/utils/file"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

const (
	defaultLogLevel       = consoleLevel
	fatalLevel      uint8 = iota
	panicLevel
	errorLevel
	warningLevel
	infoLevel
	consoleLevel
	debugLevel
)

type logOrigin struct {
	fileName string
	line     string
	fullPath string
	funcName string
}

type logCon struct {
	Message         []interface{}
	MessageStr      string
	ColorMessageStr string
	level           uint8
	origin          []*logOrigin
}

func NewLogCon(message ...interface{}) *logCon {
	return &logCon{
		Message: message,
	}
}

func logCaller(lc *logCon) (*logCon, error) {
	if lc.level == debugLevel {
		for skip := 0; ; skip++ {
			pc, file, line, ok := runtime.Caller(skip)
			if !ok {
				return lc, nil
			}
			funcInfo := runtime.FuncForPC(pc)
			if funcInfo == nil {
				return lc, errors.New("error when call runtime.FuncForPC")
			}
			lgn := &logOrigin{
				fileName: filepath.Base(file),
				line:     strconv.Itoa(line),
				fullPath: file,
				funcName: funcInfo.Name(),
			}
			lc.origin = append(lc.origin, lgn)
		}
	}
	return lc, nil
}

func (lc *logCon) setPrefix(prefix ...string) {
	preMess := make([]interface{}, len(prefix)+1)
	for k, _ := range preMess {
		if k == 0 {
			preMess[k] = "[" + prefix[k] + "] "
			continue
		}
		if k == len(preMess)-1 {
			preMess[k] = time.Now().Format("2006-01-02 15:04:05")
			continue
		}
		preMess[k] = prefix[k]
	}
	lc.Message = append(preMess[:], lc.Message...)
}

func (lc *logCon) String() string {
	return fmt.Sprintf("%s", lc.Message)
}

//添加日志队列
type LogQueue struct {
	loggers   chan *logCon
	level     uint8
	mux       *sync.Mutex
	logOutWay []int //[0]表示默认终端，[1]表示输出到文件，[2]输出到邮件； [0 1]表示双输出

	stdout *log.Logger //默认的终端输出(不区分级别)

	//文件日志
	logWriter        *log.Logger
	errLogWriter     *log.Logger
	logOutPath       string
	logOutFile       string   //日志输出文件路径
	errLogOutFile    string   //错误日志输出文件路径
	logOutFile_FD    *os.File //日志输出文件路径的文件引用
	errLogOutFile_FD *os.File //错误日志输出文件路径的文件引用
	maxLogFileSize   int64    //设置日志文件大小上线，大于这个值后重现生成日志文件

	//邮件日志
	smtpLogger *FlyWhisper.SMTPSender
}

func (lq *LogQueue) SetLevel(level int) {
	lq.level = uint8(level)
}

//设置输出位置logOutWay
func (lq *LogQueue) SetDestination(dest ...int) {
	lq.logOutWay = dest
}

//添加输出位置
func (lq *LogQueue) AddDestination(dest ...int) {
alldest:
	for _, dest := range dest {
		for _, way := range lq.logOutWay {
			if way == dest {
				continue alldest
			}
		}
		lq.logOutWay = append(lq.logOutWay, dest)
	}
}

//移除某一输出
func (lq *LogQueue) RemoveDestination(key int) {
	lenway := len(lq.logOutWay)
	for k, v := range lq.logOutWay {
		if v == key {
			if k == lenway-1 {
				lq.logOutWay = lq.logOutWay[:k]
			} else {
				lq.logOutWay = append(lq.logOutWay[:k], lq.logOutWay[k+1:]...)
			}
			return
		}
	}
}

//获取输出位置
func (lq *LogQueue) GetDestination() []int {
	return lq.logOutWay
}

//判断某一输出状态是否存在
func (lq *LogQueue) DestExist(key int) bool {
	for _, v := range lq.logOutWay {
		if v == key {
			return true
		}
	}
	return false
}

func (lq *LogQueue) FATAL(message ...interface{}) {
	lc := NewLogCon(message...)
	lc.setPrefix("FATAL")
	for _, mess := range lc.Message {
		lc.ColorMessageStr += WhiteText(RedAntiWhiteText(mess)) + " "
		lc.MessageStr += fmt.Sprintf("%s", mess) + " "
	}
	lq.intoQueue(lc, fatalLevel)
	time.Sleep(100 * time.Microsecond)
	os.Exit(1)
}

func (lq *LogQueue) ERROR(message ...interface{}) {
	lc := NewLogCon(message...)
	lc.setPrefix("ERROR")
	for _, mess := range lc.Message {
		lc.ColorMessageStr += RedText(mess) + " "
		lc.MessageStr += fmt.Sprintf("%s", mess) + " "
	}
	lq.intoQueue(lc, errorLevel)
}

func (lq *LogQueue) WARNING(message ...interface{}) {
	lc := NewLogCon(message...)
	lc.setPrefix("WARNING")
	for _, mess := range lc.Message {
		lc.ColorMessageStr += YellowText(mess) + " "
		lc.MessageStr += fmt.Sprintf("%s", mess) + " "
	}
	lq.intoQueue(lc, warningLevel)
}

func (lq *LogQueue) INFO(message ...interface{}) {
	lc := NewLogCon(message...)
	lc.setPrefix("INFO")
	for _, mess := range lc.Message {
		lc.ColorMessageStr += GreenText(mess) + " "
		lc.MessageStr += fmt.Sprintf("%s", mess) + " "
	}
	lq.intoQueue(lc, infoLevel)
}

func (lq *LogQueue) intoQueue(logCon *logCon, level uint8) {
	logCon.level = level
	if logCon.level <= lq.level {
		lq.loggers <- logCon
	}
}

func (lq *LogQueue) DEBUG(message ...interface{}) (*logCon, error) {
	lc := NewLogCon(message...)
	lc.setPrefix("DEBUG", "---\r\n")
	lc.level = debugLevel
	if lc.level <= lq.level {
		lc, err := logCaller(lc)
		if err != nil {
			log.Println("error when debugLevel called in runtime")
			return nil, err
		}
		originMess := ""
		for _, o := range lc.origin {
			originMess += "\r\n" + o.fileName + " line:" + o.line + " path:" + o.fullPath + " func:" + o.funcName + "\r\n---"
		}
		lc.Message = append(lc.Message, originMess)
		for _, mess := range lc.Message {
			lc.ColorMessageStr += BlueBoldText(mess) + " "
			lc.MessageStr += fmt.Sprintf("%s", mess) + " "
		}
		lq.loggers <- lc
		return lc, nil
	}
	return nil, errors.New("logLevel is lower than debugLevel")
}

//不会真实调用panic，需要请直接panic，但是不会记录日志文件
func (lq *LogQueue) PANIC(message ...interface{}) (*logCon, error) {
	lc := NewLogCon(message...)
	lc.setPrefix("PANIC", "---\r\n")
	lc.level = panicLevel
	if lc.level <= lq.level {
		lc, err := logCaller(lc)
		if err != nil {
			log.Println("error when debugLevel called in runtime")
			return nil, err
		}
		lc.Message = append(lc.Message, "\r\n", debug.Stack(), "\r\n---")
		for _, mess := range lc.Message {
			lc.ColorMessageStr += RedBoldText(mess) + " "
			lc.MessageStr += fmt.Sprintf("%s", mess) + " "
		}
		lq.loggers <- lc
		return lc, nil
	}
	return nil, errors.New("logLevel is lower than debugLevel")
}

func (lq *LogQueue) DoLogQueue() {
	qLevel := lq.level
	for tlog := range lq.loggers {
		if tlog.level <= qLevel {
			if lq.DestExist(0) {
				//只有在终端输出时才会显示颜色
				lq.stdout.Println(tlog.ColorMessageStr)
			}
			if lq.DestExist(1) {
				//先判断是否需要更新日志文件（日志文件size大于maxsize）
				lq.updateLogFile(lq.logOutFile)
				lq.updateErrLogFile(lq.errLogOutFile)
				fd, _ := lq.logOutFile_FD.Stat()
				//判断日志文件描述符是否关闭或者文件是否被删除
				if fd == nil || !file.IsExist(lq.logOutFile) {
					f, err := lq.getFileFd(lq.logOutFile)
					if err == nil {
						lq.SetLogOut(f)
					}
				}
				//判断错误日志文件描述符是否关闭或者文件是否被删除
				errfd, _ := lq.errLogOutFile_FD.Stat()
				if errfd == nil || !file.IsExist(lq.errLogOutFile) {
					f, err := lq.getFileFd(lq.errLogOutFile)
					if err == nil {
						lq.SetErrLogOut(f)
					}
				}
				if tlog.level <= warningLevel {
					//warning一下等级输出到错误日志
					lq.errLogWriter.Println(tlog.MessageStr)
				} else {
					lq.logWriter.Println(tlog.MessageStr)
				}
			}
		}
	}
}

//输出日志到对应位置（终端，文件）
func (lq *LogQueue) logToDestination(mess string, level uint8) {
	if lq.DestExist(0) {
		lq.stdout.Println(mess)
	}
	if lq.DestExist(1) {
		if level <= warningLevel {
			//warning一下等级输出到错误日志
			lq.errLogWriter.Println(mess)
		} else {
			lq.logWriter.Println(mess)
		}
	}

}

//可以自定义输出颜色，输出前缀
func (lq *LogQueue) Console(lc *logCon) {
	lq.intoQueue(lc, consoleLevel)
}

func NewLogger() *LogQueue {
	logger := &LogQueue{
		loggers:        make(chan *logCon, 10000),
		level:          defaultLogLevel,
		mux:            new(sync.Mutex),
		maxLogFileSize: 20 << 20, //20MB  默认日志文件大小上限，超过会自动创建新日志文件
	}
	//设置默认只输出到os.Stdout
	logger.SetDestination(0)
	logger.stdout = log.New(os.Stdout, "", 0)
	go logger.DoLogQueue()
	return logger
}
