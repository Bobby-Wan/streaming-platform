package configure

import (
	"log"

	"github.com/spf13/viper"
)

type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

var Config = viper.New()

func Init() {
	Config.AddConfigPath("..\\")
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")

	err := Config.ReadInConfig()
	if err != nil {
		log.Fatal("could not read config.yaml")
	}
	Config.MergeInConfig()
}
