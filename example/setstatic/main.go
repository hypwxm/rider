package main

import (
	"rider"
	"os"
	"path/filepath"
)

func main() {
	app := rider.New()
	app.Logger(8)
	wd, _ := os.Getwd()
	app.SetStatic(filepath.Join(wd, "src/rider/example/setStatic/public"))
	app.Listen(":5001")
}
