package logger

import (
	"testing"
	"time"
	"rider/utils"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	logger.SetDestination(0)
	logger.SetLevel(8)
	if len(logger.GetDestination()) != 1 || logger.GetDestination()[0] != 0 {
		t.Error("SetDestination error")
	}
	dir, _ := utils.GetDirName(2)
	t.Logf("%s", dir)
	logger.SetLogOutPath(dir)
	//logger.AddDestination(1)
	if len(logger.GetDestination()) != 2 || logger.GetDestination()[0] != 0 {
		t.Error("SetDestination error")
	}
	logger.INFO("xxx")
	logger.PANIC("xxx")
	logger.DEBUG("xxx")
	logger.WARNING("xxx")
	logger.FATAL("xxx")
	time.Sleep(2e9)
}