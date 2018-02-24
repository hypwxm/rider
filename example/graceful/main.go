package main

import "rider"

func main() {
	app := rider2.New()
	app.GET("/graceful", func(c rider2.Context) {

	})

	app.Graceful(":5008")
}
