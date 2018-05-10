package main

import (
	"fmt"
<<<<<<< HEAD

=======
>>>>>>> 66429bb2e8a26d25e006cfadcc518719597dfd3f
	"github.com/hypwxm/rider"
)

func main() {
	//目前不支持params和pathParams混写
	app := rider.New()
<<<<<<< HEAD
	app.GET("/id/id2/:id3", func(c rider.Context) {
=======
	app.GET("/:id/id2/:id3", func(c rider.Context) {
>>>>>>> 66429bb2e8a26d25e006cfadcc518719597dfd3f
		fmt.Println(c.Params())
		fmt.Println(c.Param("id"))
		fmt.Println(c.Param("id2"))
		c.Send(200, []byte("ok1"))
	})

	app.GET("/id/id2/(.*)", func(c rider.Context) {
		fmt.Println(c.Params())
		fmt.Println(c.PathParams())
		c.Send(200, []byte("ok2"))
	})

	//pathParams解析正则匹配到的数据
<<<<<<< HEAD
	app.GET("/reg/([a-z]{1}([a-z]{1}([a-z]{1})))/(.*)", func(c rider.Context) {
=======
	app.GET("/reg/([a-z]{1}([a-z]{1}))/(.*)", func(c rider.Context) {
>>>>>>> 66429bb2e8a26d25e006cfadcc518719597dfd3f
		fmt.Println(c.PathParams())
		c.Send(200, []byte("ok3"))
	})
	app.Listen(":5001")
}
