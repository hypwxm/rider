package main

import (
	"rider"
	//user2 "github.com/hypwxm/rider/example/router2"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/hypwxm/rider/example/router/router"
)

type pint int

func main() {

	env := flag.String("env", rider.ENV_Production, "设置app环境变量")
	cpuProfile := flag.String("cpuProfile", "", "xx")

	flag.Parse()

	rider.SetEnvMode(*env)

	//new一个rider，创建一个app
	app := rider.New()
	app.Logger(8)

	app.Kid("/super", router.Router())

	startCPUProfile(cpuProfile)
	app.Listen(":8000")
}

func startCPUProfile(cpuProfile *string) {

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not create cpu profile output file: %s",
				err)
			return
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Can not start cpu profile: %s", err)
			f.Close()
			return
		}
	}
}
