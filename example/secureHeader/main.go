package main

import (
	"net/http"
	"rider"
	"rider/modules"
)

func main() {
	app := rider2.New()
	app.USE(modules.SecureHeader())
	app.GET("/secure", func(context rider2.Context) {
		context.Send(http.StatusOK, []byte("...."))
	})
	app.Listen(":5002")
}
