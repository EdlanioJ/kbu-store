package routes

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/factory"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func storeRoutes(route fiber.Router, db *gorm.DB, tc time.Duration) {
	getByID := factory.GetStoreByIDHandler(db, tc)
	getByOwner := factory.GetStoreByOwnerHandler(db, tc)
	fetch := factory.FetchStoreHandler(db, tc)
	fetchByStatus := factory.FetchStoreByStatusHandler(db, tc)
	fetchByType := factory.FetchStoreByTypeHandler(db, tc)
	fetchByOwner := factory.FetchStoreByOwnerHandler(db, tc)
	fetchByTags := factory.FetchStoreByTagsHandler(db, tc)
	fetchByLocation := factory.FetchByLocationAndStatusHandler(db, tc)
	block := factory.BlockStoreHandler(db, tc)
	active := factory.ActiveStoreHandler(db, tc)
	disable := factory.DisableStoreHandler(db, tc)
	create := factory.CreateStoreHandler(db, tc)
	delete := factory.DeleteStoreHandler(db, tc)
	update := factory.UpdateStoreHandler(db, tc)

	storeRoutes := route.Group("/stores")

	storeRoutes.Get("/", fetch.Handler)
	storeRoutes.Get("/status/:status", fetchByStatus.Handler)
	storeRoutes.Get("/location/:location/status/:status", fetchByLocation.Handler)
	storeRoutes.Get("/category/:category", fetchByType.Handler)
	storeRoutes.Get("/owner/:owner", fetchByOwner.Handler)
	storeRoutes.Get("/tags/", fetchByTags.Handler)
	storeRoutes.Get("/:id", getByID.Handler)
	storeRoutes.Get("/:id/owner/:owner", getByOwner.Handler)

	storeRoutes.Post("/", create.Handler)
	storeRoutes.Delete("/:id", delete.Handler)
	storeRoutes.Patch("/:id", update.Handler)

	storeRoutes.Options("/:id/block", block.Handler)
	storeRoutes.Options("/:id/active", active.Handler)
	storeRoutes.Options("/:id/disable", disable.Handler)
}
