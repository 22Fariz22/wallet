package wallet

import (
	"github.com/labstack/echo/v4"
)

type Handlers interface {
	Operation() echo.HandlerFunc
	Display() echo.HandlerFunc
	CreateWallet() echo.HandlerFunc
}
