package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/config"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/http"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/kafka"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/EdlanioJ/kbu-store/app/usecases"
	"github.com/spf13/cobra"
)

var httpPort int

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "start http server",
	Run: func(*cobra.Command, []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}

		database := repository.GORMConnection(cfg)

		httpServer := http.NewHttpServer()

		httpServer.Port = cfg.Port
		if httpPort != 0 {
			httpServer.Port = httpPort
		}

		tc := time.Duration(cfg.Timeout) * time.Second
		kafkaProducer := kafka.NewKafkaProducer(cfg)
		storeRepo := gorm.NewStoreRepository(database)
		accountRepo := gorm.NewAccountRepository(database)
		categoryRepo := gorm.NewCategoryRepository(database)

		httpServer.StoreUsecase = usecases.NewStoreUsecase(
			storeRepo,
			accountRepo,
			categoryRepo,
			kafkaProducer,
			tc,
		)

		httpServer.Serve()
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().IntVarP(&httpPort, "port", "p", 0, "http server port")
}
