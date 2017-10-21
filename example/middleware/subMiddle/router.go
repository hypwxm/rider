package subMiddle

import (
	"rider"
	"fmt"
)

func Router() *rider.Router {
	router := rider.NewRouter()
	router.AddMiddleware(
		func(c *rider.Context) {
			fmt.Println("sub1")
			c.Next()
		},
	)
	router.GET("/subMiddle", &rider.Router{
		Handler: func(c *rider.Context) {
			fmt.Println("sub2")
			c.Send([]byte("ok"))
		},
		Middleware: []rider.HandlerFunc{
			func(c *rider.Context) {
				fmt.Println("insub2")
				c.Next()
			},
		},
	})


	subrouter := &rider.Router{
		Middleware: []rider.HandlerFunc{
			func(c *rider.Context) {
				fmt.Println("mid in mid")
				c.Next()
			},
		},
	}

	router.GET("/mid", subrouter)
	subrouter.GET("midd", &rider.Router{
		Handler: func(c *rider.Context) {
			fmt.Println("mide in mid in mid")
		},
	})

	return router
}
