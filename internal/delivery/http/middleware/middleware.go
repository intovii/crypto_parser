package middleware

import (
	"go.uber.org/zap"
	"CryptoParser/config"
	"CryptoParser/internal/repository/postgres"
)

type Middleware struct {
	cfg   *config.ConfigModel
	repo  *postgres.Repository
	log   *zap.Logger
	roles map[string]int
}

func NewMiddleware(cfg *config.ConfigModel, log *zap.Logger, repository *postgres.Repository) *Middleware {
	return &Middleware{
		cfg:  cfg,
		log:  log,
		repo: repository,
	}
}
