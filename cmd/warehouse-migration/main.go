package main

import (
	"github.com/texazcowboy/warehouse/cmd/warehouse-migration/application"
	// register drivers.
	_ "github.com/lib/pq"
)

// note: applies all scripts
func main() {
	var app application.App
	app.Initialize()
	app.Run()
}
