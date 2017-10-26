package main

import (
	"rider"
	"rider/modules"
	"net/http"
)

func main() {
	app := rider3.New()
	app.AddMiddleware(modules.SecureHeader())
	app.GET("/secure", &rider3.Router{
		Handler: func(context *rider3.Context) {
			context.Send(http.StatusOK, []byte("...."))
		},
	})
	app.Listen(":5002")
}
