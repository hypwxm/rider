package router

import (
	"rider"
	"fmt"
)

func Router() *rider.Router {
	router := rider.NewRouter()
	router.GET("/", &rider.Router{
		Handler: func (c *rider.Context) {
			fmt.Println(c.GetLocals("locals"))
			fmt.Println(c.GetLocals("locals2"))
			c.Send([]byte("ok"))
		},
	})
	router.GET("/xx", &rider.Router{
		Handler: func (c *rider.Context) {
			fmt.Println(c.GetLocals("locals"))
			fmt.Println(c.GetLocals("locals2"))
			c.Send([]byte("ok2"))
		},
	})
	return router
}
