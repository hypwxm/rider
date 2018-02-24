package main

import (
	"fmt"
	"rider"
)

func main() {
	app := rider2.New()
	app.GET("/", func(c rider2.Context) {
		auth := c.HeaderValue("Authorization")
		//b, err := base64.URLEncoding.DecodeString(strings.Split(auth, " ")[1])
		fmt.Println(auth)

		c.SetHeader(rider2.HeaderWWWAuthenticate, "Digest realm='a and b' qop='adad' nonce='adad'")
		c.Send(401, []byte(""))
	})
	app.GET("/auth", func(c rider2.Context) {

		fmt.Println(c.Header())
	})
	app.Listen(":5001")
}
