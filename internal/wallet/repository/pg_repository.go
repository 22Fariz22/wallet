package repository

import (
	"context"

	"github.com/22Fariz22/wallet/internal/wallet"
	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type walletRepo struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewWalletRepository(db *sqlx.DB, logger logger.Logger) wallet.Repository {
	return &walletRepo{db: db, logger: logger}
}

func (r walletRepo) Display(context context.Context, walletID uuid.UUID) (int64, error) {
	return 0, nil
}

func (r *walletRepo) Deposit(context context.Context, walletID uuid.UUID, amount int64) error {
	return nil
}

func (r *walletRepo) Withdraw(context context.Context, walletID uuid.UUID, amount int64) error {
	return nil
}
