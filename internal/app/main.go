package app

import (
	"context"
	"fmt"

	"github.com/linspacestrom/go-project/internal/auth"
	"github.com/linspacestrom/go-project/internal/config"
	"github.com/linspacestrom/go-project/internal/event"
	"github.com/linspacestrom/go-project/internal/repository"
	"github.com/linspacestrom/go-project/internal/server"
	authService "github.com/linspacestrom/go-project/internal/service/auth"
	"github.com/linspacestrom/go-project/internal/service/outbox"
	platformService "github.com/linspacestrom/go-project/internal/service/platform"
	authHandler "github.com/linspacestrom/go-project/internal/transport/http/handlers/auth"
	infoHandler "github.com/linspacestrom/go-project/internal/transport/http/handlers/info"
	platformHandler "github.com/linspacestrom/go-project/internal/transport/http/handlers/platform"
	swaggerHandler "github.com/linspacestrom/go-project/internal/transport/http/handlers/swagger"
	"go.uber.org/zap"
)

type Repository interface {
	Close()
	Ping(ctx context.Context) error
}

type Producer interface {
	Close() error
}

type App struct {
	log         *zap.Logger
	api         *server.Server
	cfg         *config.Config
	repo        Repository
	outboxSvc   *outbox.Service
	outboxCtx   context.Context
	outboxStop  context.CancelFunc
	kafkaWriter Producer
}

func New(log *zap.Logger, cfg *config.Config) (*App, error) {
	repo, trManager, err := repository.New(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	tokenManager := auth.NewManager(cfg.Auth.Secret, cfg.Auth.TokenTTL)

	authSvc := authService.NewService(tokenManager, repo, cfg.Auth.RefreshTokenTTL)
	authHndl := authHandler.NewHandler(authSvc)

	platformSvc := platformService.NewService(repo, trManager)
	platformHndl := platformHandler.NewHandler(platformSvc)

	infoHndl := infoHandler.NewHandler(repo)
	swaggerHndl := swaggerHandler.NewHandler()

	publicHandlers := []server.Handler{
		infoHndl,
		authHndl,
		swaggerHndl,
	}

	protectedHandlers := []server.Handler{platformHndl}

	api := server.New(log, cfg.Server, cfg.Auth.Secret, publicHandlers, protectedHandlers)

	producer := event.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	outboxSvc := outbox.NewService(log, repo, producer, cfg.Outbox.DispatchInterval, cfg.Outbox.BatchSize)
	outboxCtx, outboxStop := context.WithCancel(context.Background())

	return &App{
		log:         log,
		api:         api,
		cfg:         cfg,
		repo:        repo,
		outboxSvc:   outboxSvc,
		outboxCtx:   outboxCtx,
		outboxStop:  outboxStop,
		kafkaWriter: producer,
	}, nil
}

func (a *App) Run() {
	defer func() {
		if r := recover(); r != nil {
			a.log.Error("application panicked", zap.Any("panic", r))
			a.Stop()
		}
	}()

	go a.outboxSvc.Run(a.outboxCtx)
	a.api.MustRun()
}

func (a *App) Stop() {
	a.log.Info("closing HTTP server")
	if err := a.api.Close(); err != nil {
		a.log.Error("failed to close HTTP server", zap.Error(err))
	}

	a.log.Info("stopping outbox dispatcher")
	a.outboxStop()

	a.log.Info("closing kafka producer")
	if err := a.kafkaWriter.Close(); err != nil {
		a.log.Error("failed to close kafka producer", zap.Error(err))
	}

	a.log.Info("closing database connection pool")
	a.repo.Close()

	a.log.Info("application stopped gracefully")
}
