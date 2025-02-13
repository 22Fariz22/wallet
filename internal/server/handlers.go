package server

import (
	walletHTTP "github.com/22Fariz22/wallet/internal/wallet/delivery/http"
	walletRepository "github.com/22Fariz22/wallet/internal/wallet/repository"
	walletUseCase "github.com/22Fariz22/wallet/internal/wallet/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) MapHandlers(e *echo.Echo) error {

	// Init repositories
	walletRepo := walletRepository.NewWalletRepository(s.db, s.logger)

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

	v1 := e.Group(s.cfg.Middleware.MiddlewareAPIVersion)

	walletGroup := v1.Group("/")

	walletHTTP.MapWalletRoutes(walletGroup, walletHandler)

	return nil
}
