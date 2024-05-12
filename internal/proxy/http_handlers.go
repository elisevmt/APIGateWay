package proxy

import "github.com/gofiber/fiber/v2"

type HTTPHandlers interface {
	Proxy() fiber.Handler
}
