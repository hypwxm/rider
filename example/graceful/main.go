package main

import (
	"rider"
)

func main() {
	app := rider.New()
	app.GET("/graceful", func(c rider.Context) {

	})

	app.Graceful(":5008")
}
