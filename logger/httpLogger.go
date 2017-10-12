package logger

import (
	"time"
	"strconv"
)

func HttpLogger(lq *LogQueue, method string, path string, statusCode int, duration time.Duration, exMess ...string)  {
	var codeColor string
	codeStr := strconv.Itoa(statusCode)
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
	lq.Console(GreenText("[HTTP] "), WhiteText(method), WhiteText(path), codeColor, WhiteText(duration), GreenText("=>"), BlueText(exMess))
}