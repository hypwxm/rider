package main

import (
	"rider"
	//user2 "rider/example/router2"
	"fmt"
	"flag"
	"rider/example/router/router"
)

type pint int


func main() {

	env := flag.String("env", rider.ENV_Development, "设置app环境变量")
	flag.Parse()

	rider.SetEnvMode(*env)

	//new一个rider，创建一个app
	app := rider.New()

	app.ANY("/super", router.Router())

	app.Listen(":8000")
}
