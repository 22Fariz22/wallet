package server

import (
	"log"
	"os"

	walletHTTP "github.com/22Fariz22/wallet/internal/wallet/delivery/http"
	walletRepository "github.com/22Fariz22/wallet/internal/wallet/repository"
	walletUseCase "github.com/22Fariz22/wallet/internal/wallet/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	s.logger.Info("Registering routes...")

	// открываем файл для подсчета количества записей после нагрузочного тестирования
	file, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(file)

	// Init repositories
	walletRepo := walletRepository.NewWalletRepository(s.db, s.logger, s.redisClient)

	// Init useCases
	walletUC := walletUseCase.NewWalletUseCase(s.cfg, walletRepo, s.redisClient, s.logger)

	// Init handlers
	walletHandler := walletHTTP.NewWalletHandler(s.cfg, walletUC, s.logger)

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         s.cfg.Middleware.MiddlewareStackSize,
		DisablePrintStack: s.cfg.Middleware.MiddlewareDisablePrintStack,
		DisableStackAll:   s.cfg.Middleware.MiddlewareDisableStackAll,
	}))
	e.Use(middleware.RequestID())

	s.logger.Debug("API Version:", s.cfg.Middleware.MiddlewareAPIVersion)
	v1 := e.Group(s.cfg.Middleware.MiddlewareAPIVersion)

	walletGroup := v1.Group("")

	walletHTTP.MapWalletRoutes(walletGroup, walletHandler)

	return nil
}
