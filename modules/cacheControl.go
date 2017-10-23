package modules

import (
	"rider"
	"time"
	"net/http"
	"strconv"
)

/*
type cacheControl struct {
	CacheControl string `header:"Cache-Control"`  //max-age=100
	Expires string `header:"Expires"`  //Sun, 22 Oct 2017 14:30:43 GMT
}
*/

//request里面也有cache-control，会影响实际请求中是否使用缓存还是向服务器请求

func CacheControl(d int64) rider.HandlerFunc {
	dStr := strconv.FormatInt(int64(d), 10)
	return func(c *rider.Context) {
		c.SetHeader("Cache-Control", "max-age=" + dStr)
		c.SetHeader("Expires", time.Now().UTC().Add(time.Duration(d)*time.Second).Format(http.TimeFormat))
		c.Next()
	}
}