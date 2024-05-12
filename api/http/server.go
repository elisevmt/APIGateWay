package http

import (
	"APIGateWay/api"
	"APIGateWay/config"
	"APIGateWay/constants"
	"APIGateWay/internal"
	"APIGateWay/internal/proxy"
	proxy_delivery_http "APIGateWay/internal/proxy/delivery/http"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recoverMDW "github.com/gofiber/fiber/v2/middleware/recover"
)

type httpServer struct {
	fiber *fiber.App
	cfg   *config.Config
}

func NewHTTPServer(cfg *config.Config) api.Server {
	return &httpServer{
		cfg: cfg,
	}
}

func (h *httpServer) Init() error {
	h.fiber = fiber.New(fiber.Config{
		Immutable:               true,
		AppName:                 "APIGateWay",
		EnableTrustedProxyCheck: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if logicErr, ok := err.(constants.LogicError); ok {
				if logicErr.Code == -1 {
					return ctx.Status(200).SendString(logicErr.Message)
				}
				return ctx.Status(int(logicErr.Code)).SendString(logicErr.Message)
			}
			return ctx.Status(500).SendString("Unknown error")
		},
	})
	return nil
}

func (h *httpServer) MapHandlers(app *internal.App) error {
	h.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	h.fiber.Use(recoverMDW.New())
	h.fiber.Use(logger.New())
	h.fiber.Get("/version", func(ctx *fiber.Ctx) error {
		err := ctx.JSON(h.cfg.Server.Version)
		if err != nil {
			return err
		}
		return nil
	})
	transactionHandlers := proxy_delivery_http.NewProxyHTTPHandlers(app.UC["proxy"].(proxy.UC), app.GuzzleLogger)
	proxyGroup := h.fiber.Group("proxy")
	proxy_delivery_http.MapTransactionRoutes(proxyGroup, transactionHandlers)
	return nil
}

func (h *httpServer) Run() error {
	fmt.Printf("LISTENING %s:%s\n", h.cfg.Server.Host, h.cfg.Server.Port)
	err := h.fiber.Listen(h.cfg.Server.Host + ":" + h.cfg.Server.Port)
	return err
}
