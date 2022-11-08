package config

import (
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	VkApi VkApiConfig
}

type AppConfig struct {
	Port        string
	LogLevel    string
	TokenSecret string
}

type DBConfig struct {
	FSConf string
}

type VkApiConfig struct {
	BotToken string
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

	config.App.TokenSecret = os.Getenv("TOKEN_SECRET")
	if config.App.TokenSecret == "" {
		return nil, errors.New("Not specified TOKEN_SECRET")
	}

	//db parse

	config.DB.FSConf = os.Getenv("FS_CONF")
	if config.DB.FSConf == "" {
		return nil, errors.New("Not specified FS_CONF")
	}

	//vk api

	config.VkApi.BotToken = os.Getenv("BOT_TOKEN")
	if config.VkApi.BotToken == "" {
		return nil, errors.New("Not specified BOT_TOKEN")
	}

	return &config, err

}
