package middleware

import "github.com/gofiber/fiber/v2"

func NotFound(message ...string) fiber.Handler {
	msg := fiber.ErrNotFound.Message
	if len(message) > 0 {
		msg = message[0]
	}

	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": msg,
		})
	}
}
