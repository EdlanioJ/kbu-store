package cmd

import (
	"time"

	"github.com/EdlanioJ/kbu-store/app/config"
	"github.com/EdlanioJ/kbu-store/app/factory"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/kafka"
	"github.com/EdlanioJ/kbu-store/app/infrastructure/repository"
	"github.com/spf13/cobra"
)

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start kafka consumer",
	Run: func(*cobra.Command, []string) {
		cfg, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}

		database := repository.GORMConnection(cfg)

		tc := time.Duration(cfg.Timeout) * time.Second

		kafkaCosumer := kafka.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.GroupID)
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
