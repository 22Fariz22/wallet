package wallet

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Deposit(context context.Context, walletID uuid.UUID, amount int64) error
	Withdraw(context context.Context, walletID uuid.UUID, amount int64) error
	Display(context context.Context, walletID uuid.UUID) (int64, error)
}
