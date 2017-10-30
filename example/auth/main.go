package main

import (
	"rider"
	"fmt"
)

func main() {
	app := rider.New()
	app.GET("/", func(c rider.Context) {
		auth := c.HeaderValue("Authorization")
		//b, err := base64.URLEncoding.DecodeString(strings.Split(auth, " ")[1])
		fmt.Println(auth)

		c.SetHeader(rider.HeaderWWWAuthenticate, "Digest realm='a and b' qop='adad' nonce='adad'")
		c.Send(401, []byte(""))
	})
	app.GET("/auth", func(c rider.Context) {

		fmt.Println(c.Header())
	})
	app.Listen(":5001")
}
