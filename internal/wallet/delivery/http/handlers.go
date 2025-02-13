package http

import (
	"net/http"

	"github.com/22Fariz22/wallet/config"
	"github.com/22Fariz22/wallet/internal/wallet"
	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/labstack/echo/v4"
)

type walletHandlers struct {
	cfg           *config.Config
	walletUsecase wallet.Usecase
	logger        logger.Logger
}

func NewWalletHandler(
	cfg *config.Config,
	walletUsecase wallet.Usecase,
	logger logger.Logger,
) wallet.Handlers {
	return &walletHandlers{cfg: cfg, walletUsecase: walletUsecase, logger: logger}
}

func (h walletHandlers) Display() echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

func (h walletHandlers) Deposit() echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}

func (h walletHandlers) Withdraw() echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.NoContent(http.StatusOK)
	}
}
