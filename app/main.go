package main

import (
	"fmt"
	"time"

	"github.com/EdlanioJ/kbu-store/app/config"
	_ "github.com/EdlanioJ/kbu-store/app/docs"
	"github.com/EdlanioJ/kbu-store/app/routes"
	"github.com/EdlanioJ/kbu-store/app/utils"
	"github.com/sirupsen/logrus"
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
	config, err := config.LoadConfig("./")
	if err != nil {
		logrus.Fatal("can not import config file", err)
	}

	timeoutContext := time.Duration(config.TimeoutContext) * time.Second
	db := utils.ConnectDB()

	app := routes.New(db, timeoutContext)

	logrus.Fatal(app.Listen(fmt.Sprintf(":%d", config.Port)))
}
