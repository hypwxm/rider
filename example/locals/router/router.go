package router

import (
	"rider"
	"fmt"
)

func Router() *rider3.Router {
	router := rider3.NewRouter()
	router.GET("/", &rider3.Router{
		Handler: func (c *rider3.Context) {
			fmt.Println(c.GetLocals("locals"))
			fmt.Println(c.GetLocals("locals2"))
			c.Send([]byte("ok"))
		},
	})
	router.GET("/xx", &rider3.Router{
		Handler: func (c *rider3.Context) {
			fmt.Println(c.GetLocals("locals"))
			fmt.Println(c.GetLocals("locals2"))

			c.Send([]byte("ok2"))
		},
	})
	return router
}
