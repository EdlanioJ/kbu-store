package main

import (
	"time"

	"github.com/EdlanioJ/kbu-store/application/config"
	"github.com/EdlanioJ/kbu-store/application/http"
	_ "github.com/EdlanioJ/kbu-store/application/http/docs"
	"github.com/EdlanioJ/kbu-store/infra/db"
	"gorm.io/gorm"
)

func main() {
	var database *gorm.DB
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	tc := time.Duration(config.TimeoutContext) * time.Second

	if config.Env == "test" {
		database = db.GORMConnection(config.DnsTest, config.Env)
	} else {
		database = db.GORMConnection(config.Dns, config.Env)
	}

	http.StartServer(database, tc, config.Port)
}
