package main

import (
	"bpzh-api/internal/client/vk_api"
	"bpzh-api/internal/config"
	"bpzh-api/internal/controller"
	"bpzh-api/internal/logic"
	"bpzh-api/internal/repo/db"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed read Config")
	}

	dbConnectionPool, err := db.CreateFSConnections(&cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize DB")
	}
	defer dbConnectionPool.Close()

	vkApiClient := vk_api.NewClient(&cfg.VkApi)

	repo := db.NewRepo(dbConnectionPool)          // model работает с БД и прочими источниками данных
	lgc := logic.NewLogic(cfg, repo, vkApiClient) // logic знает, что делать с model
	api := controller.NewApp(cfg, lgc)            // api использует logic для обработки запросов

	api.StartServe()
}
