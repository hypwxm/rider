package main

import "rider"

func main() {
	app := rider.New()
	logger := app.Logger(8)
	logger.INFO("adad")
	app.GET("/", func(c rider.Context) {
			c.Redirect(301, "http://localhost:3000")
	})
	app.Listen(":5001")
}
