package config

import (
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	Port     string
	LogLevel string
}

type DBConfig struct {
	FSConf string
}

func ReadConfig() (*Config, error) {

	var config Config
	var err error

	//app parse

	config.App.Port = os.Getenv("SERVER_PORT")
	if config.App.Port == "" {
		config.App.Port = "80"
	}

	config.App.LogLevel = os.Getenv("LOG_LEVEL")
	if config.App.LogLevel == "" {
		config.App.LogLevel = "debug"
	}

	//db parse

	config.DB.FSConf = os.Getenv("FS_CONF")
	if config.DB.FSConf == "" {
		return nil, errors.New("Not specified FS_CONF")
	}

	return &config, err

}
