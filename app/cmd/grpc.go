package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/config"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/kafka"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository/gorm"
	"github.com/EdlanioJ/kbu-store/app/usecases"
	"github.com/spf13/cobra"
)

var port int

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "start grpc server",
	Run: func(*cobra.Command, []string) {
		cfg, err := config.LoadConfig(".")
		if err != nil {
			panic(err)
		}

		database := repository.GORMConnection(cfg)

		grpcServer := grpc.NewGrpcServer()

		grpcServer.Port = cfg.Grpc.Port
		if port != 0 {
			grpcServer.Port = port
		}

		tc := time.Duration(cfg.Timeout) * time.Second
		kafkaProducer := kafka.NewKafkaProducer(cfg)
		storeRepo := gorm.NewStoreRepository(database)
		accountRepo := gorm.NewAccountRepository(database)
		categoryRepo := gorm.NewCategoryRepository(database)

		grpcServer.MetricPort = cfg.Grpc.MetricPort
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

		grpcServer.StoreUsecase = storeUsecase

		grpcServer.Serve()
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().IntVarP(&port, "port", "p", 0, "grpc server port")
}
