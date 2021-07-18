package http

import (
	"fmt"
	"log"
	"time"

	"github.com/EdlanioJ/kbu-store/application/factory"
	_ "github.com/EdlanioJ/kbu-store/application/http/docs"
	"github.com/EdlanioJ/kbu-store/application/http/handler"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"gorm.io/gorm"
)

// @title KBU Store API
// @version 2.0.0
// @description This is a sample swagger for KBU Store
// @termsOfService http://swagger.io/terms/
// @contact.name Edlâneo Manuel
// @contact.email edlanioj@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
func StartServer(database *gorm.DB, tc time.Duration, port int) {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(requestid.New())
	v1 := app.Group("/api/v1")

	v1.Get("/docs/*", swagger.Handler)
	storeUsecase := factory.StoreUsecase(database, tc)
	tu := factory.TagUsecase(database, tc)
	cu := factory.CategoryUsecase(database, tc)

	storeRoutes := v1.Group("/stores")
	storeHandler := handler.NewStoreHandler(storeUsecase)

	storeRoutes.Post("/", storeHandler.Store)
	storeRoutes.Get("/", storeHandler.Index)
	storeRoutes.Get("/:id", storeHandler.Get)
	storeRoutes.Patch("/:id", storeHandler.Update)
	storeRoutes.Patch("/:id/activate", storeHandler.Activate)
	storeRoutes.Patch("/:id/block", storeHandler.Block)
	storeRoutes.Patch("/:id/disable", storeHandler.Disable)
	storeRoutes.Delete("/:id", storeHandler.Delete)

	handler.NewTagRoutes(v1, tu)
	handler.NewCategoryRoutes(v1, cu)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
