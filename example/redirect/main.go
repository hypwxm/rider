package main

import "rider"

func main() {
	app := rider2.New()
	logger := app.Logger(8)
	logger.INFO("adad")
	app.GET("/", func(c rider2.Context) {
		c.Redirect(301, "http://localhost:3000")
	})
	app.Listen(":5001")
}
