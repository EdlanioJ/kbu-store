package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/config"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/http"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/jaeger"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/kafka"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/EdlanioJ/kbu-store/app/usecases"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
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

		tracer, closer, err := jaeger.InitJaeger(cfg)
		if err != nil {
			log.Fatal("cannot create tracer ", err)
		}

		opentracing.SetGlobalTracer(tracer)
		defer closer.Close()

		database := repository.GORMConnection(cfg)

		httpServer := http.NewHttpServer()

		httpServer.Port = cfg.Port
		if httpPort != 0 {
			httpServer.Port = httpPort
		}

		kafkaProducer := kafka.NewKafkaProducer(cfg)
		defer kafkaProducer.Close()

		tc := time.Duration(cfg.Timeout) * time.Second
		storeRepo := gorm.NewStoreRepository(database)
		accountRepo := gorm.NewAccountRepository(database)
		categoryRepo := gorm.NewCategoryRepository(database)

		storeUsecase := usecases.NewStoreUsecase(
			storeRepo,
			accountRepo,
			categoryRepo,
			kafkaProducer,
			tc,
		)

		storeUsecase.NewStoreTopic = cfg.Kafka.NewStoreTopic
		storeUsecase.UpdateStoreTopic = cfg.Kafka.UpdateStoreTopic
		storeUsecase.DeleteStoreTopic = cfg.Kafka.DeleteStoreTopic

		httpServer.StoreUsecase = storeUsecase
		httpServer.Validate = validator.New()

		httpServer.Serve()
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().IntVarP(&httpPort, "port", "p", 0, "http server port")
}
