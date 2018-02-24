package main

import "rider"

func main() {
	app := rider2.New()
	app.GET("/hijack", func(c rider2.Context) {
		c.Hijack()
		c.SetHeader("a", "b")
		c.Send(200, []byte("ok"))
		c.SendJson(200, "adad")
		c.SendFile("path")
	})
}
