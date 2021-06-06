package routes

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/factory"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func tagRoutes(route fiber.Router, db *gorm.DB, tc time.Duration) {
	fetchTagsHandler := factory.FetchTagsHandler(db, tc)
	fetchTagsByCategoryHandler := factory.FetchTagsByCategoryHandler(db, tc)

	tagRoutes := route.Group("/tags")

	tagRoutes.Get("/", fetchTagsHandler.Handler)
	tagRoutes.Get("/category/:category", fetchTagsByCategoryHandler.Handler)

}
