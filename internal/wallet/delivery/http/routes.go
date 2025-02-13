package http

import (
	"github.com/22Fariz22/wallet/internal/wallet"
	"github.com/labstack/echo/v4"
)

func MapWalletRoutes(walletGroup *echo.Group, h wallet.Handlers) {
	walletGroup.GET("/wallets/:uiid", h.Display())
	walletGroup.POST("/wallet", h.Deposit())
	walletGroup.POST("/wallet", h.Withdraw())
}
