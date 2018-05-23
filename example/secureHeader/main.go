package main

import (
	"net/http"
	"rider"

	"github.com/hypwxm/rider/modules"
)

func main() {
	app := rider.New()
	app.USE(modules.SecureHeader())
	app.GET("/secure", func(context rider.Context) {
		context.Send(http.StatusOK, []byte("...."))
	})
	app.Listen(":5002")
}
