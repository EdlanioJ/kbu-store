package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Timeout      int      `mapstructure:"TIMEOUT"`
	Port         int      `mapstructure:"PORT"`
	GrpcPort     int      `mapstructure:"GRPC_PORT"`
	MetricPort   int      `mapstructure:"METRIC_PORT"`
	Env          string   `mapstructure:"ENV"`
	Dns          string   `mapstructure:"PG_DNS"`
	DnsTest      string   `mapstructure:"PG_DNS_TEST"`
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
}

func LoadConfig(path ...string) (config Config, err error) {
	defaultPath := "."

	if len(path) > 0 {
		defaultPath = path[0]
	}

	viper.AddConfigPath(defaultPath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 3333)
	viper.SetDefault("GRPC_PORT", 50051)
	viper.SetDefault("METRIC_PORT", 3330)
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
