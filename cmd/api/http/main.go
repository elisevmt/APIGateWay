package main

import (
	"APIGateWay/api/http"
	"APIGateWay/config"
	"APIGateWay/internal"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		panic(err)
	}
	app := internal.NewApp(cfg)
	httpServer := http.NewHTTPServer(cfg)
	err = app.Init()
	if err != nil {
		panic(err)
	}
	err = httpServer.Init()
	if err != nil {
		panic(err)
	}
	err = httpServer.MapHandlers(app)
	if err != nil {
		panic(err)
	}
	err = httpServer.Run()
	if err != nil {
		panic(err)
	}
}
