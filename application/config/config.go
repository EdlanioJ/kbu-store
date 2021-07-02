package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TimeoutContext int    `mapstructure:"TIMEOUT_CONTEXT"`
	Port           int    `mapstructure:"PORT"`
	Env            string `mapstructure:"ENV"`
	Dns            string `mapstructure:"PG_DNS"`
	DnsTest        string `mapstructure:"PG_DNS_TEST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 3333)
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
