package logger

import (
	"time"
	"strconv"
	"strings"
)

func HttpLogger(lq *LogQueue, method string, path string, statusCode int, duration time.Duration, exMess ...string) {
	var codeColor string
	codeStr := strconv.Itoa(statusCode)
	lc := NewLogCon("[HTTP] ", method, path, codeStr, duration, "=>", exMess)

	switch {
	case statusCode >= 500:
		codeColor = RedBg(YellowText(codeStr))
	case statusCode >= 400:
		codeColor = RedAntiWhiteText(codeStr)
	case statusCode >= 300:
		codeColor = YellowAntiWhiteText(codeStr)
	case statusCode >= 200:
		codeColor = BlueAntiWhiteText(codeStr)
	case statusCode >= 100:
		codeColor = WhiteAntiWhiteText(codeStr)
	}
	lc.ColorMessageStr = strings.Join([]string{GreenText("[HTTP] "), WhiteText(method), WhiteText(path), codeColor, WhiteText(duration), GreenText("=>"), BlueText(exMess)}, " ")
	messArr := []string{"[HTTP] ", method, path, codeStr, duration.String(), "=>", "["}
	messArr = append(messArr, exMess...)
	messArr = append(messArr, "]")
	lc.MessageStr = strings.Join(messArr, " ")
	lq.Console(lc)
}
