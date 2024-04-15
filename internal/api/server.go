package api

import (
	"context"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"time"

	"github.com/uptrace/bunrouter"

	"github.com/KseniiaSalmina/Car-catalog/internal/config"
	"github.com/KseniiaSalmina/Car-catalog/internal/database_manager"
	"github.com/KseniiaSalmina/Car-catalog/internal/info_receiver"
	"github.com/KseniiaSalmina/Car-catalog/internal/logger"
	"github.com/KseniiaSalmina/Car-catalog/internal/models"
)

type Server struct {
	dbManager    dbManager
	infoReceiver infoReceiver
	httpServer   *http.Server
	logger       *logger.Logger
}

type dbManager interface {
	DeleteCar(ctx context.Context, regNum string) error
	ChangeCar(ctx context.Context, car models.Car) error
	AddCars(ctx context.Context, cars []models.Car) error
	GetCars(ctx context.Context, filters models.Car, yearFilterMode string, orderByMode string, limit, offset int) (*models.CarsPage, error)
}

type infoReceiver interface {
	GetCarInfo(ctx context.Context, regNum string) (*models.Car, error)
}

func NewServer(cfg config.Server, dbManager *database_manager.PostgresManager, receiver *info_receiver.Receiver, logger *logger.Logger) *Server {
	s := &Server{
		dbManager:    dbManager,
		infoReceiver: receiver,
	}

	router := bunrouter.New(bunrouter.WithMiddleware(s.middlewareLog)).Compat()
	router.GET("/cars", s.GetCars)
	router.POST("/cars", s.PostCars)
	router.PATCH("/cars/:regNum", s.PatchCar)
	router.DELETE("/cars/:regNum", s.DeleteCar)

	swagHandler := httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json"))
	router.GET("/swagger/*path", swagHandler)

	s.httpServer = &http.Server{
		Addr:         cfg.Listen,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return s
}

func (s *Server) Run() {
	s.logger.Logger.Info("server started")

	go func() {
		err := s.httpServer.ListenAndServe()
		s.logger.Logger.Info(fmt.Sprintf("http server stopped: %s", err.Error()))
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
