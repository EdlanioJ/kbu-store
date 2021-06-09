package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

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
	var dns string
	dbConnention := os.Getenv("PG_DNS")
	dbConnectionTest := os.Getenv("PG_DNS_TEST")
	port := os.Getenv("PORT")
	tc, _ := strconv.Atoi(os.Getenv("TIMEOUT_CONTEXT"))
	env := os.Getenv("ENV")
	migration := os.Getenv("AUTO_MIGRATE_DB")

	if env != "test" {
		dns = dbConnention
	} else {
		dns = dbConnectionTest
	}

	timeoutContext := time.Duration(tc) * time.Second
	db := utils.ConnectDB(env, dns, migration)

	app := routes.New(db, timeoutContext)

	logrus.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
