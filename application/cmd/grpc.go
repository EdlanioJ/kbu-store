package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/application/config"
	"github.com/EdlanioJ/kbu-store/application/grpc"
	"github.com/EdlanioJ/kbu-store/infra/db"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var port int

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "start grpc server",
	Run: func(*cobra.Command, []string) {
		var database *gorm.DB
		config, err := config.LoadConfig(".")
		if err != nil {
			panic(err)
		}

		if config.Env == "test" {
			database = db.GORMConnection(config.DnsTest, config.Env)
		} else {
			database = db.GORMConnection(config.Dns, config.Env)
		}

		tc := time.Duration(config.TimeoutContext) * time.Second
		grpc.StartServer(database, tc, port)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
	grpcCmd.Flags().IntVarP(&port, "port", "p", 50051, "grpc server port")
}