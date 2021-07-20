package http

import (
	"fmt"
	"log"

	_ "github.com/EdlanioJ/kbu-store/application/http/docs"
	"github.com/EdlanioJ/kbu-store/application/http/handler"
	"github.com/EdlanioJ/kbu-store/application/http/middleware"
	"github.com/EdlanioJ/kbu-store/domain"
	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
)

type httpServer struct {
	Port            int
	StoreUsecase    domain.StoreUsecase
	CategoryUsecase domain.CategoryUsecase
}

func NewHttpServer() *httpServer {
	return &httpServer{}
}

// @title KBU Store API
// @version 2.0.0
// @description This is a sample swagger for KBU Store
// @termsOfService http://swagger.io/terms/
// @contact.name Edlâneo Manuel
// @contact.email edlanioj@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
func (s *httpServer) Serve() {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
	})

	prometheus := fiberprometheus.New("kbu-store")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(requestid.New())

	v1 := app.Group("/api/v1")

	v1.Get("/docs/*", swagger.Handler)

	storeHandler := handler.NewStoreHandler(s.StoreUsecase)
	storeRoutes := v1.Group("/stores")

	storeRoutes.Post("/", storeHandler.Store)
	storeRoutes.Get("/", storeHandler.Index)
	storeRoutes.Get("/:id", storeHandler.Get)
	storeRoutes.Patch("/:id", storeHandler.Update)
	storeRoutes.Patch("/:id/activate", storeHandler.Activate)
	storeRoutes.Patch("/:id/block", storeHandler.Block)
	storeRoutes.Patch("/:id/disable", storeHandler.Disable)
	storeRoutes.Delete("/:id", storeHandler.Delete)

	handler.NewCategoryRoutes(v1, s.CategoryUsecase)

	app.Use(middleware.NotFound())

	log.Fatal(app.Listen(fmt.Sprintf(":%d", s.Port)))
}
