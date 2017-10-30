package main

import (
	"rider"
	"rider/modules"
	"net/http"
)

func main() {
	app := rider.New()
	app.USE(modules.SecureHeader())
	app.GET("/secure", func(context rider.Context) {
			context.Send(http.StatusOK, []byte("...."))
	})
	app.Listen(":5002")
}
