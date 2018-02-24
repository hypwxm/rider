package router

import (
	"fmt"
	"rider"
)

func Router() *rider2.Router {
	router := rider2.NewRouter()
	router.GET("/", func(c rider2.Context) {
		fmt.Println(c.GetLocals("locals"))
		fmt.Println(c.GetLocals("locals2"))
		c.Send(200, []byte("ok"))
	})
	router.GET("/xx", func(c rider2.Context) {
		fmt.Println(c.GetLocals("locals"))
		fmt.Println(c.GetLocals("locals2"))
		c.Send(200, []byte("ok2"))
	})
	return router
}
