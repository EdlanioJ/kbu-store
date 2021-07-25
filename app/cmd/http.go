package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/config"
	"github.com/EdlanioJ/kbu-store/app/factory"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/http"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository"
	"github.com/spf13/cobra"
)

var httpPort int

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "start http server",
	Run: func(*cobra.Command, []string) {
		config, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}

		database := repository.GORMConnection(config)

		tc := time.Duration(config.Timeout) * time.Second

		httpServer := http.NewHttpServer()
		httpServer.Port = config.Port

		if httpPort != 0 {
			httpServer.Port = httpPort
		}
		httpServer.StoreUsecase = factory.StoreUsecase(database, tc, config.Kafka.Brokers)

		httpServer.Serve()
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().IntVarP(&httpPort, "port", "p", 0, "http server port")
}
