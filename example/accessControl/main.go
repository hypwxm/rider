// 跨域信息设置
package main

import "github.com/hypwxm/rider"

func main() {
	app := rider.New()
	app.SetAccessCtl(func(c rider.Context) *rider.AccessControl {
		return &rider.AccessControl{
			AccessControlAllowCredentials: "true",
			AccessControlAllowOrigin:      "http://localhost:8080",
			AccessControlAllowHeaders:     "content-type",
			AccessControlAllowMethods:     "get",
		}
	})
}
