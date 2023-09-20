// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"tm-backend-trainee-impl-clean-template/config"
	v1 "tm-backend-trainee-impl-clean-template/internal/controller/http/v1"
	"tm-backend-trainee-impl-clean-template/internal/usecase"
	"tm-backend-trainee-impl-clean-template/internal/usecase/repo"
	"tm-backend-trainee-impl-clean-template/pkg/httpserver"
	"tm-backend-trainee-impl-clean-template/pkg/logger"
	"tm-backend-trainee-impl-clean-template/pkg/postgres"
	"tm-backend-trainee-impl-clean-template/pkg/validator"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	statisticsUseCase := usecase.New(
		repo.New(pg),
	)

	validator.RegisterCustomValidators(l)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, statisticsUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
