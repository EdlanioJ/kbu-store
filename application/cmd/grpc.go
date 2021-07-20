package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/application/config"
	"github.com/EdlanioJ/kbu-store/application/factory"
	"github.com/EdlanioJ/kbu-store/application/grpc"
	"github.com/EdlanioJ/kbu-store/infra/db"
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

		database := db.GORMConnection(config.Dns, config.Env)
		if config.Env == "test" {
			database = db.GORMConnection(config.DnsTest, config.Env)
		}

		tc := time.Duration(config.Timeout) * time.Second
		grpcServer := grpc.NewGrpcServer()

		grpcServer.Port = config.GrpcPort
		if port != 0 {
			grpcServer.Port = port
		}

		grpcServer.MetricPort = config.MetricPort
		grpcServer.CategoryUsecase = factory.CategoryUsecase(database, tc)
		grpcServer.StoreUsecase = factory.StoreUsecase(database, tc, config.KafkaBrokers)

		grpcServer.Serve()
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().IntVarP(&port, "port", "p", 0, "grpc server port")
}
