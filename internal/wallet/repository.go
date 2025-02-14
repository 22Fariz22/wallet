package wallet

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error
	Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error
	Display(ctx context.Context, walletID uuid.UUID) (int64, error)
	CreateWallet(ctx context.Context) (uuid.UUID, error)
}
