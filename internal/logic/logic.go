package logic

import (
	"bpzh-api/internal/client/vk_api"
	"bpzh-api/internal/config"
	"bpzh-api/internal/repo/db"
)

// Logic содержит все для доступа к данным
type Logic struct {
	config          *config.Config
	repo            *db.Repo
	vkApiBpzhClient *vk_api.Client
}

func NewLogic(config *config.Config, repo *db.Repo, vkApiBpzhClient *vk_api.Client) *Logic {
	return &Logic{
		config:          config,
		repo:            repo,
		vkApiBpzhClient: vkApiBpzhClient,
	}
}
