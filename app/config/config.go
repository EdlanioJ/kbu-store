package config

import (
	"github.com/spf13/viper"
)

type Kafka struct {
	GroupID     string   `mapstructure:"GROUP_ID"`
	Brokers     []string `mapstructure:"BROKERS"`
	GroupTopics []string `mapstructure:"GROUP_TOPICS"`
}

type Grpc struct {
	Port       int `mapstructure:"PORT"`
	MetricPort int `mapstructure:"METRIC_PORT"`
}

type Config struct {
	Timeout int    `mapstructure:"TIMEOUT"`
	Port    int    `mapstructure:"PORT"`
	Env     string `mapstructure:"ENV"`
	Dns     string `mapstructure:"PG_DNS"`
	DnsTest string `mapstructure:"PG_DNS_TEST"`
	Kafka   Kafka  `mapstructure:"KAFKA"`
	Grpc    Grpc   `mapstructure:"GRPC"`
}

func LoadConfig(path ...string) (cfg *Config, err error) {
	defaultPath := "."

	if len(path) > 0 {
		defaultPath = path[0]
	}

	viper.AddConfigPath(defaultPath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 3333)
	viper.SetDefault("GRPC.PORT", 50051)
	viper.SetDefault("GRPC.METRIC_PORT", 3330)
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}
