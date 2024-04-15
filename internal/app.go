package app

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
	"time"

	"github.com/KseniiaSalmina/Car-catalog/internal/api"
	"github.com/KseniiaSalmina/Car-catalog/internal/config"
	"github.com/KseniiaSalmina/Car-catalog/internal/database_manager"
	"github.com/KseniiaSalmina/Car-catalog/internal/info_receiver"
	"github.com/KseniiaSalmina/Car-catalog/internal/logger"
	"github.com/KseniiaSalmina/Car-catalog/internal/storage/postgres"
)

type Application struct {
	cfg          config.Application
	logger       *logrus.Logger
	db           *postgres.DB
	dbManager    *database_manager.PostgresManager
	infoReceiver *info_receiver.Receiver
	server       *api.Server
	closeCtx     context.Context
	closeCtxFunc context.CancelFunc
}

func NewApplication(cfg config.Application) (*Application, error) {
	app := Application{
		cfg: cfg,
	}

	if err := app.bootstrap(); err != nil {
		return nil, err
	}

	app.readyToShutdown()

	return &app, nil
}

func (a *Application) bootstrap() error {
	//init logger
	if err := a.initLogger(); err != nil {
		return fmt.Errorf("failed to bootstrap app: %w", err)
	}

	//init dependencies
	if err := a.initDatabase(); err != nil {
		return fmt.Errorf("failed to bootstrap app: %w", err)
	}

	//init services
	a.initDbManager()
	a.initInfoReceiver()

	//init controllers
	a.initServer()

	return nil
}

func (a *Application) initLogger() error {
	l, err := logger.NewLogger(a.cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}

	a.logger = l
	return nil
}

func (a *Application) initDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := postgres.NewDB(ctx, a.cfg.Postgres)
	if err != nil {
		a.logger.WithError(err).Error("failed to init database")
		return err
	}

	a.db = db
	a.logger.Info("successful connection to database")
	return nil
}

func (a *Application) initDbManager() {
	dbManager := database_manager.NewPostgresManager(a.db)
	a.dbManager = dbManager
	a.logger.Debug("successful init db manager")
}

func (a *Application) initInfoReceiver() {
	infoReceiver := info_receiver.NewReceiver(a.cfg.Receiver)
	a.infoReceiver = infoReceiver
	a.logger.Debug("successful init info receiver")
}

func (a *Application) initServer() {
	s := api.NewServer(a.cfg.Server, a.dbManager, a.infoReceiver, a.logger)
	a.server = s
	a.logger.Debug("successful init server")
}

func (a *Application) Run() {
	defer a.stop()
	a.logger.Debug("application started")

	a.server.Run()

	<-a.closeCtx.Done()
	a.closeCtxFunc()
}

func (a *Application) stop() {
	if err := a.server.Shutdown(); err != nil {
		a.logger.Errorf("incorrect closing of server: %s", err.Error())
	} else {
		a.logger.Info("server closed")
	}

	a.db.Close()
	a.logger.Info("database closed")
}

func (a *Application) readyToShutdown() {
	ctx, closeCtx := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	a.closeCtx = ctx
	a.closeCtxFunc = closeCtx
}
