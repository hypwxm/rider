package main

import "rider"

func main() {
	app := rider.New()
	app.GET("/hijack", func(c rider.Context) {
		c.Hijack()
		c.SetHeader("a", "b")
		c.Send(200, []byte("ok"))
		c.SendJson(200, "adad")
		c.SendFile("path")
	})
}
