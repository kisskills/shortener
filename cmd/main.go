package main

import "shortener/internal/application"

func main() {
	var confPath string = "./configs/service.yml"

	app := application.Application{}
	app.Build(confPath)

	app.Run()
}
