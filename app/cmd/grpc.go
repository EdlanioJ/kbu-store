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
		cfg, err := config.LoadConfig(".")
		if err != nil {
			panic(err)
		}

		database := repository.GORMConnection(cfg)

		tc := time.Duration(cfg.Timeout) * time.Second
		grpcServer := grpc.NewGrpcServer()

		grpcServer.Port = cfg.Grpc.Port
		if port != 0 {
			grpcServer.Port = port
		}

		grpcServer.MetricPort = cfg.Grpc.MetricPort
		grpcServer.StoreUsecase = factory.StoreUsecase(database, tc, cfg)

		grpcServer.Serve()
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().IntVarP(&port, "port", "p", 0, "grpc server port")
}
