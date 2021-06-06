package config

import "github.com/spf13/viper"

type Config struct {
	PGDns          string `mapstructure:"PG_DNS"`
	TimeoutContext int    `mapstructure:"TIMEOUT_CONTEXT"`
	Port           int    `mapstructure:"PORT"`
	AutoMigration  bool   `mapstructure:"AUTO_MIGRATE_DB"`
	Env            string `mapstructure:"ENV"`
	PGDnsTest      string `mapstructure:"PG_DNS_TEST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 3333)
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
