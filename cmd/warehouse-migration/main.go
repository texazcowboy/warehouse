package main

import (
	"github.com/texazcowboy/warehouse/cmd/warehouse-migration/application"
)

// note: applies all scripts
func main() {
	var app application.App
	app.Initialize()
	app.Run()
}
