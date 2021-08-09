package http

import (
	"fmt"
	"log"

	"github.com/EdlanioJ/kbu-store/app/domain"
	_ "github.com/EdlanioJ/kbu-store/app/infrastructure/http/docs"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/http/handler"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/http/middleware"
	fiberprometheus "github.com/ansrivas/fiberprometheus/v2"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
)

type httpServer struct {
	Port         int
	StoreUsecase domain.StoreUsecase
	Validate     *validator.Validate
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
	app := fiber.New()

	prometheus := fiberprometheus.New("kbu-store")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(requestid.New())

	v1 := app.Group("/api/v1")

	v1.Get("/docs/*", swagger.Handler)

	s.routes(v1)
	app.Use(middleware.NotFound())

	log.Fatal(app.Listen(fmt.Sprintf(":%d", s.Port)))
}

func (s *httpServer) routes(route fiber.Router) {
	storeHandler := handler.NewStoreHandler(s.StoreUsecase, s.Validate)
	storeRoutes := route.Group("/stores")

	storeRoutes.Post("/", storeHandler.Store)
	storeRoutes.Get("/", storeHandler.Index)
	storeRoutes.Get("/:id", storeHandler.Get)
	storeRoutes.Patch("/:id", storeHandler.Update)
	storeRoutes.Patch("/:id/activate", storeHandler.Activate)
	storeRoutes.Patch("/:id/block", storeHandler.Block)
	storeRoutes.Patch("/:id/disable", storeHandler.Disable)
	storeRoutes.Delete("/:id", storeHandler.Delete)
}
