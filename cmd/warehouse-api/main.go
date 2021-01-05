package main

import (
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/application"
)

func main() {
	var app application.App
	app.Initialize()
	app.Run()
}
