package main

import (
	"auth-service/internal/app"
	"flag"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "", "config file dir")
	flag.Parse()

	application := app.NewApplication(path)
	application.Run()
}
