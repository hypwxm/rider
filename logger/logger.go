package logger

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"errors"
	"strconv"
	"time"
	"fmt"
	"sync"
	"runtime/debug"
)

const (
	defaultLogLevel = errorLevel
	fatalLevel uint8 = iota
	panicLevel
	errorLevel
	warningLevel
	infoLevel
	debugLevel
	consoleLevel
)


type Logger interface {

}


type logOrigin struct {
	fileName string
	line     string
	fullPath string
	funcName string
}

type logCon struct {
	Message []interface{}
	level uint8
	origin *logOrigin
}

func NewLogCon(message ...interface{}) *logCon {
	return &logCon{
		Message: message,
	}
}

func logCaller(lc *logCon) (*logCon, error){
	if lc.level == debugLevel {
		pc, file, line, ok := runtime.Caller(0)
		if !ok {
			return nil, errors.New("error when call runtime.Caller")
		}
		funcInfo := runtime.FuncForPC(pc)
		if funcInfo == nil {
			return nil, errors.New("error when call runtime.FuncForPC")
		}
		lgn := &logOrigin{
			fileName: filepath.Base(file),
			line: strconv.Itoa(line),
			fullPath: file,
			funcName: funcInfo.Name(),
		}
		lc.origin = lgn
	}
	return lc, nil
}

func (lc *logCon) setPrefix(prefix ...string) {
	preMess := make([]interface{}, len(prefix) + 1)
	for k, _ := range preMess{
		if k == 0 {
			preMess[k] = "[" + prefix[k] + "] "
			continue
		}
		if k == len(preMess) - 1 {
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
	loggers chan *logCon
	level uint8
	mux *sync.Mutex
}

func (lq *LogQueue) SetLevel(level int) {
	lq.level = uint8(level)
}

func (lq *LogQueue) FATAL(message ...interface{}) {
	lc := NewLogCon(message...)
	lc.setPrefix("FATAL")
	lq.intoQueue(lc, fatalLevel)
	time.Sleep(100 * time.Microsecond)
}

func (lq *LogQueue) ERROR(message ...interface{}) {
	lc := NewLogCon(message...)
	lc.setPrefix("ERROR")
	lq.intoQueue(lc, errorLevel)
}

func (lq *LogQueue) WARNING(message ...interface{}) {
	lc := NewLogCon(message...)
	lc.setPrefix("WARNING")
	lq.intoQueue(lc, warningLevel)
}

func (lq *LogQueue) INFO(message ...interface{}) {
	lc := NewLogCon(message...)
	lc.setPrefix("INFO")
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
		lc.Message = append(lc.Message, "\r\n", debug.Stack())
		lq.loggers <- lc
		return lc, nil
	}
	return nil, errors.New("logLevel is lower than debugLevel")
}

func (lq *LogQueue) DoLogQueue() {
	qLevel := lq.level
	logOut := log.New(os.Stdout, "", 0)
	for tlog := range lq.loggers {
		if tlog.level <= qLevel {
			if tlog.level == debugLevel {
				originMess := "\r\n" + tlog.origin.fileName + " line:" + tlog.origin.line + " path:" + tlog.origin.fullPath + " func:" + tlog.origin.funcName + "\r\n---"
				tlog.Message = append(tlog.Message, originMess)
			}
			messConbin := ""
			for _, mess := range tlog.Message {
				switch tlog.level {
				case fatalLevel:
					messConbin += BlueText(GreenBg(mess))
				case panicLevel:
					messConbin += RedBoldText(mess)
				case errorLevel:
					messConbin += RedText(mess)
				case warningLevel:
					messConbin += YellowText(mess)
				case infoLevel:
					messConbin += GreenText(mess)
				case debugLevel:
					messConbin += BlueBoldText(mess)
				default:
					messConbin += fmt.Sprintf("%s", mess)
				}
				messConbin += " "
			}
			logOut.Println(messConbin)
		}
	}
}

func (lq *LogQueue) Console(message ...interface{}) {
	lc := NewLogCon(message...)
	lq.intoQueue(lc, consoleLevel)
}

func NewLogger() *LogQueue {
	logger := &LogQueue{
		loggers: make(chan *logCon, 10000),
		level: defaultLogLevel,
		mux: new(sync.Mutex),
	}
	go logger.DoLogQueue()
	return logger
}
