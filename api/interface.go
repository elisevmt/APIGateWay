package api

import "APIGateWay/internal"

type Server interface {
	Init() error
	MapHandlers(app *internal.App) error
	Run() error
}
