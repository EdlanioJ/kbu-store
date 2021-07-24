package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/config"
	"github.com/EdlanioJ/kbu-store/app/factory"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/grpc"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository"
	"github.com/spf13/cobra"
)

var port int

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "start grpc server",
	Run: func(*cobra.Command, []string) {
		config, err := config.LoadConfig(".")
		if err != nil {
			panic(err)
		}

		database := repository.GORMConnection(config.Dns, config.Env)
		if config.Env == "test" {
			database = repository.GORMConnection(config.DnsTest, config.Env)
		}

		tc := time.Duration(config.Timeout) * time.Second
		grpcServer := grpc.NewGrpcServer()

		grpcServer.Port = config.Grpc.Port
		if port != 0 {
			grpcServer.Port = port
		}

		grpcServer.MetricPort = config.Grpc.MetricPort
		grpcServer.StoreUsecase = factory.StoreUsecase(database, tc, config.Kafka.Brokers)

		grpcServer.Serve()
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().IntVarP(&port, "port", "p", 0, "grpc server port")
}
