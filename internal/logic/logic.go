package logic

import (
	"bpzh-api/internal/config"
	"bpzh-api/internal/repo/db"
)

// Logic содержит все для доступа к данным
type Logic struct {
	config *config.Config
	repo   *db.Repo
}

func NewLogic(config *config.Config, repo *db.Repo) *Logic {
	return &Logic{
		config: config,
		repo:   repo,
	}
}
