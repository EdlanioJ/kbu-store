package config

import (
	"github.com/spf13/viper"
)

type Kafka struct {
	GroupID             string   `mapstructure:"GROUP_ID"`
	Brokers             []string `mapstructure:"BROKERS"`
	CreateCategoryTopic string   `mapstructure:"CREATE_CATEGORY_TOPIC"`
	UpdateCategoryTopic string   `mapstructure:"UPDATE_CATEGORY_TOPIC"`
	NewStoreTopic       string   `mapstructure:"NEW_STORE_TOPIC"`
	UpdateStoreTopic    string   `mapstructure:"UPDATE_STORE_TOPIC"`
	DeleteStoreTopic    string   `mapstructure:"DELETE_STORE_TOPIC"`
}

type Grpc struct {
	Port       int `mapstructure:"PORT"`
	MetricPort int `mapstructure:"METRIC_PORT"`
}

type Jaeger struct {
	Host        string `mapstructure:"HOST"`
	ServiceName string `mapstructure:"SERVICE_NAME"`
	LogSpans    bool   `mapstructure:"LOG_SPANS"`
}

type Config struct {
	Timeout int    `mapstructure:"TIMEOUT"`
	Port    int    `mapstructure:"PORT"`
	Env     string `mapstructure:"ENV"`
	Dns     string `mapstructure:"DB_CONNECTION"`
	DnsTest string `mapstructure:"DB_CONNECTION_TEST"`
	Kafka   Kafka  `mapstructure:"KAFKA"`
	Grpc    Grpc   `mapstructure:"GRPC"`
	Jaeger  Jaeger `mapstructure:"JEAGER"`
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
