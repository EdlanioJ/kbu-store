package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/application/config"
	"github.com/EdlanioJ/kbu-store/application/factory"
	"github.com/EdlanioJ/kbu-store/application/http"
	"github.com/EdlanioJ/kbu-store/infra/db"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "start http server",
	Run: func(*cobra.Command, []string) {
		var database *gorm.DB
		config, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}

		if config.Env == "test" {
			database = db.GORMConnection(config.DnsTest, config.Env)
		} else {
			database = db.GORMConnection(config.Dns, config.Env)
		}

		tc := time.Duration(config.Timeout) * time.Second

		httpServer := http.NewHttpServer()

		httpServer.StoreUsecase = factory.StoreUsecase(database, tc, config.KafkaBrokers)
		httpServer.CategoryUsecase = factory.CategoryUsecase(database, tc)
		httpServer.Port = config.Port

		httpServer.Serve()
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
