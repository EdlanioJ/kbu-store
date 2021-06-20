package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/EdlanioJ/kbu-store/app/config"
	_ "github.com/EdlanioJ/kbu-store/app/docs"
	"github.com/EdlanioJ/kbu-store/app/factory"
	"github.com/EdlanioJ/kbu-store/app/utils"
	categoryRoute "github.com/EdlanioJ/kbu-store/category/deliver/http"
	storeRoute "github.com/EdlanioJ/kbu-store/store/deliver/http"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
)

// @title KBU Store API
// @version 1.0.0
// @description This is a sample swagger for KBU Store
// @termsOfService http://swagger.io/terms/
// @contact.name Edlâneo Manuel
// @contact.email edlanioj@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3333
// @BasePath /api/v1
func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	tc := time.Duration(config.TimeoutContext) * time.Second
	db := utils.ConnectDB()

	app := fiber.New()

	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${ua} ${latency} ${ip} ${locals:requestid} ${method} ${path}​\n​",
		Output: os.Stderr,
	}))

	v1 := app.Group("/api/v1")
	v1.Get("/docs/*", swagger.Handler)
	su := factory.StoreUsecase(db, tc)
	tu := factory.TagUsecase(db, tc)
	cu := factory.CategoryUsecase(db, tc)

	storeRoute.NewStoreRoute(v1, su)
	storeRoute.NewTagRoutes(v1, tu)
	categoryRoute.NewCategoryRoutes(v1, cu)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.Port)))
}
