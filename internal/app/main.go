package app

import (
	"fmt"

	"github.com/linspacestrom/go-project/internal/config"
	"github.com/linspacestrom/go-project/internal/repository"
	"github.com/linspacestrom/go-project/internal/server"
	"go.uber.org/zap"
)

type Repository interface {
	Close()
}

type App struct {
	log  *zap.Logger
	api  *server.Server
	cfg  *config.Config
	repo Repository
}

func New(log *zap.Logger, cfg *config.Config) (*App, error) {
	repo, trManager, err := repository.New(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	_ = trManager
	api := server.New(log, cfg.Server)

	return &App{
		log:  log,
		api:  api,
		cfg:  cfg,
		repo: repo,
	}, nil
}

func (a *App) Run() {
	defer func() {
		if r := recover(); r != nil {
			a.log.Error("application panicked", zap.Any("panic", r))
			a.Stop()
		}
	}()

	a.api.MustRun()
}

func (a *App) Stop() {
	a.log.Info("closing HTTP server")
	if err := a.api.Close(); err != nil {
		a.log.Error("failed to close HTTP server", zap.Error(err))
	}

	a.log.Info("closing database connection pool")
	a.repo.Close()

	a.log.Info("application stopped gracefully")
}
