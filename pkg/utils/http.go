package utils

import (
	"github.com/gofiber/fiber/v2"
)

func ReadRequest(c *fiber.Ctx, request interface{}) error {
	if err := c.BodyParser(request); err != nil {
		return err
	}

	return validate.StructCtx(c.Context(), request)
}
