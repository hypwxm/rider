package main

import (
	"rider"
	"rider/modules"
	"net/http"
)

func main() {
	app := rider.New()
	app.AddMiddleware(modules.SecureHeader())
	app.GET("/secure", &rider.Router{
		Handler: func(context *rider.Context) {
			context.Send(http.StatusOK, []byte("...."))
		},
	})
	app.Listen(":5002")
}
