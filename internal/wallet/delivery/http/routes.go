package http

import (
	"github.com/22Fariz22/wallet/internal/wallet"
	"github.com/labstack/echo/v4"
)

func MapWalletRoutes(walletGroup *echo.Group, h wallet.Handlers) {
	walletGroup.GET("/wallets/:uuid", h.Display())
	walletGroup.POST("/wallet", h.Operation())
	walletGroup.POST("/new", h.CreateWallet())
}
