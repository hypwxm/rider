package main

import (
	"rider"
	//user2 "rider/example/router2"
	"flag"
	"rider/example/router/router"
)

type pint int


func main() {

	env := flag.String("env", rider3.ENV_Production, "设置app环境变量")
	flag.Parse()

	rider3.SetEnvMode(*env)

	//new一个rider，创建一个app
	app := rider3.New()

	app.ANY("/super", router.Router())

	app.Listen(":8000")
}
