package main

import (
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/entrypoint"
)

func main() {
	var app entrypoint.App
	app.Initialize()
	app.Run()
}
