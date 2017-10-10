package main

import (
	"rider"
	"os"
	"path/filepath"
)

func main() {
	app := rider.New()
	wd, _ := os.Getwd()
	app.SetStatic(filepath.Join(wd, "src/rider/example/setStatic/public"))
	app.Listen(":5000")
}
