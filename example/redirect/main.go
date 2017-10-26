package main

import "rider"

func main() {
	app := rider3.New()
	logger := app.Logger(8)
	logger.INFO("adad")
	app.GET("/", &rider3.Router{
		Handler: func(c *rider3.Context) {
			c.SetStatusCode(302)
			c.Redirect(301, "http://localhost:3000")
		},
	})
	app.Listen(":5001")
}
