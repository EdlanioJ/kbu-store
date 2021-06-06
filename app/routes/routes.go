package routes

import (
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"gorm.io/gorm"
)

func New(db *gorm.DB, tc time.Duration) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(logger.New())
	v1 := app.Group("/api/v1")

	v1.Get("/docs/*", swagger.Handler)

	storeRoutes(v1, db, tc)
	categoryRoutes(v1, db, tc)
	tagRoutes(v1, db, tc)

	return app
}
