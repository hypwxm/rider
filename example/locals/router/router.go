package router

import (
	"fmt"
	"rider"
)

func Router() *rider.Router {
	router := rider.NewRouter()
	router.GET("/", func(c rider.Context) {
		fmt.Println(c.GetLocals("locals"))
		fmt.Println(c.GetLocals("locals2"))
		c.Send(200, []byte("ok"))
	})
	router.GET("/xx", func(c rider.Context) {
		fmt.Println(c.GetLocals("locals"))
		fmt.Println(c.GetLocals("locals2"))
		c.Send(200, []byte("ok2"))
	})
	return router
}
