package wallet

import (
	"github.com/labstack/echo/v4"
)

type Handlers interface {
	Deposit() echo.HandlerFunc
	Withdraw() echo.HandlerFunc
	Display() echo.HandlerFunc
}
