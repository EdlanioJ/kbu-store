package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/application/config"
	"github.com/EdlanioJ/kbu-store/application/factory"
	"github.com/EdlanioJ/kbu-store/infra/db"
	"github.com/EdlanioJ/kbu-store/infra/kafka"
	"github.com/spf13/cobra"
)

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start kafka consumer",
	Run: func(*cobra.Command, []string) {
		config, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}

		database := db.GORMConnection(config.Dns, config.Env)
		if config.Env == "test" {
			database = db.GORMConnection(config.DnsTest, config.Env)
		}

		tc := time.Duration(config.Timeout) * time.Second

		kafkaCosumer := kafka.NewKafkaConsumer(config.Kafka.Brokers, config.Kafka.GroupID)
		kafkaCosumer.CategoryUsecase = factory.CategoryUsecase(database, tc)

		kafkaCosumer.Consume()
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kafkaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kafkaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
