package application

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"shortener/internal/adapters/storage/inmemory"
	"shortener/internal/adapters/storage/postgres"
	"shortener/internal/cases"
	"shortener/internal/config"
	"shortener/internal/transport/grpc"
	"syscall"
)

type Application struct {
	cancel  context.CancelFunc
	log     *zap.SugaredLogger
	cfg     *config.Config
	storage cases.Storage
	server  *grpc.Server
}

func (a *Application) Build(configPath string) {
	var err error

	a.log = a.initLogger()

	a.cfg, err = config.New(configPath)
	if err != nil {
		a.log.Fatal("init config")
	}

	if a.cfg.DatabaseMode() == "postgres" {
		a.storage = a.buildPostgresStorage()
	} else {
		a.storage = a.buildInMemoryStorage()
	}

	svc := a.buildService(a.storage)

	a.server = a.buildServer(svc)

}

func (a *Application) Run() {
	a.log.Info("application started")
	defer a.log.Info("application stopped")

	var ctx context.Context

	ctx, a.cancel = context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		a.Stop()
	}()

	a.server.Run(ctx)
}

func (a *Application) Stop() {
	a.storage.Close()
	a.cancel()
	_ = a.log.Sync()
}

func (a *Application) initLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		a.log.Fatal(err)
	}

	return logger.Sugar()
}

func (a *Application) buildPostgresStorage() *postgres.PGStorage {
	st, err := postgres.NewPGStorage(a.log, a.cfg.PostgresDSN())
	if err != nil {
		a.log.Fatal(err)
	}

	return st
}
func (a *Application) buildInMemoryStorage() *inmemory.InMemoryStorage {
	st, err := inmemory.NewInMemoryStorage(a.log)
	if err != nil {
		a.log.Fatal(err)
	}

	return st
}

func (a *Application) buildService(storage cases.Storage) *cases.UrlService {
	svc, err := cases.NewUrlService(a.log, storage)
	if err != nil {
		a.log.Fatal(err)
	}

	return svc
}

func (a *Application) buildServer(svc *cases.UrlService) *grpc.Server {
	srv, err := grpc.NewServer(a.log, svc, a.cfg.HttpPort(), a.cfg.GrpcPort())
	if err != nil {
		a.log.Fatal(err)
	}

	return srv
}
