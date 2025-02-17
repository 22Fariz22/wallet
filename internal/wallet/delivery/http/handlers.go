package http

import (
	"net/http"

	"github.com/22Fariz22/wallet/config"
	"github.com/22Fariz22/wallet/internal/wallet"
	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/22Fariz22/wallet/pkg/utils"
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
		h.logger.Info("Display handler called")

		ctx := c.Request().Context()

		uuidStr := c.Param("uuid")

		// Проверяем UUID
		walletUUID, err := utils.ValidateUUID(uuidStr)
		if err != nil {
			h.logger.Warnf("Invalid UUID: %s", uuidStr)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":   "invalid UUID format",
				"message": "UUID must be in valid format (e.g., 550e8400-e29b-41d4-a716-446655440000)",
			})
		}

		amount, err := h.walletUsecase.Display(ctx, walletUUID)
		if err != nil {
			h.logger.Errorf("Failed to fetch balance for wallet %s: %v", walletUUID, err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":   "internal server error",
				"message": "please retry query",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Balance retrieved successfully",
			"amount":  amount,
		})
	}
}

type WalletTransactionRequest struct {
	WalletID      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

func (h walletHandlers) Operation() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Info("Operation handler called")

		ctx := c.Request().Context()

		// Парсим JSON
		var req WalletTransactionRequest
		if err := c.Bind(&req); err != nil {
			h.logger.Warn("Invalid request body")
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":   "invalid request body",
				"message": "Check JSON structure",
			})
		}

		// Используем функцию из utils
		walletUUID, err := utils.ValidateUUID(req.WalletID)
		if err != nil {
			h.logger.Warnf("Invalid UUID: %s", req.WalletID)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":   "invalid UUID format",
				"message": "UUID must be in valid format (e.g., 550e8400-e29b-41d4-a716-446655440000)",
			})
		}

		switch req.OperationType {
		case "DEPOSIT":
			err = h.walletUsecase.Deposit(ctx, walletUUID, req.Amount)
		case "WITHDRAW":
			err = h.walletUsecase.Withdraw(ctx, walletUUID, req.Amount)
		default:
			h.logger.Warnf("Invalid operation type: %s", req.OperationType)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":   "invalid operation type",
				"message": "operationType must be 'DEPOSIT' or 'WITHDRAW'",
			})
		}

		if err != nil {
			h.logger.Errorf("Operation failed: %s, walletID: %s, amount: %d, error: %v",
				req.OperationType, walletUUID, req.Amount, err)

			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":   "operation failed",
				"message": "please try again later",
			})
		}

		h.logger.Infof("Operation successful: %s, walletID: %s, amount: %d",
			req.OperationType, walletUUID, req.Amount)
		return c.NoContent(http.StatusOK)
	}
}

func (h walletHandlers) CreateWallet() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Info("CreateWallet handler called")

		ctx := c.Request().Context()

		// Вызываем usecase для создания кошелька
		walletID, err := h.walletUsecase.CreateWallet(ctx)
		if err != nil {
			h.logger.Errorf("Failed to create wallet: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":   "wallet creation failed",
				"message": "please try again later",
			})
		}

		h.logger.Infof("Wallet created successfully: %s", walletID)
		return c.JSON(http.StatusCreated, map[string]string{
			"wallet_id": walletID.String(),
			"message":   "Wallet created successfully",
		})
	}
}
