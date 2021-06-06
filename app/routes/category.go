package routes

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/factory"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func categoryRoutes(route fiber.Router, db *gorm.DB, tc time.Duration) {
	getByID := factory.GetCategoryByIDHandler(db, tc)
	getByStatus := factory.GetCategoryByStatusHandler(db, tc)
	fetch := factory.FetchCategoryHandler(db, tc)
	fetchByStatus := factory.FetchCategoryByStatusHandler(db, tc)
	create := factory.CreateCategoryHandler(db, tc)

	categoryRoutes := route.Group("/categories")

	categoryRoutes.Get("/", fetch.Handler)
	categoryRoutes.Get("/status/:status", fetchByStatus.Handler)
	categoryRoutes.Get("/:id", getByID.Handler)
	categoryRoutes.Get("/:id/:status", getByStatus.Handler)
	categoryRoutes.Post("/", create.Handler)
}
