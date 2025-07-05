package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Db     *Database
		Server *Server
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DbName   string
		SslMode  string
		TimeZone string
	}

	Server struct {
		Port int
	}
)

var (
	once   sync.Once
	config *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		err := viper.Unmarshal(&config)
		if err != nil {
			panic(err)
		}
	})

	return config
}
