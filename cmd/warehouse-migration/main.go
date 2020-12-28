package main

import (
	"github.com/texazcowboy/warehouse/cmd/warehouse-migration/entrypoint"
)

// note: applies all scripts
func main() {
	var app entrypoint.App
	app.Initialize()
	app.Run()
}
