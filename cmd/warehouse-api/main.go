package main

import (
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/application"
	// register drivers.
	_ "github.com/lib/pq"
)

func main() {
	var app application.App
	app.Initialize()
	app.Run()
}
