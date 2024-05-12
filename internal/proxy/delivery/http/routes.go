package proxy_delivery_http

import (
	"APIGateWay/internal/proxy"
	"github.com/gofiber/fiber/v2"
)

func MapTransactionRoutes(router fiber.Router, h proxy.HTTPHandlers) {
	router.All("/:service_id/*", h.Proxy())
}
