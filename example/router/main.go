package main

import (
	"rider"
	//user2 "rider/example/router2"
	"flag"
	"rider/example/router/router"
)

type pint int


func main() {

	env := flag.String("env", rider.ENV_Production, "设置app环境变量")
	flag.Parse()

	rider.SetEnvMode(*env)

	//new一个rider，创建一个app
	app := rider.New()

	app.Kid("/super", router.Router())

	app.Listen(":8000")
}
