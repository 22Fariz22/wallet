package wallet

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// Мок для Usecase
type MockWalletUsecase struct {
	mock.Mock
}

func (m *MockWalletUsecase) Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error {
	args := m.Called(ctx, walletID, amount)
	return args.Error(0)
}

func (m *MockWalletUsecase) Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error {
	args := m.Called(ctx, walletID, amount)
	return args.Error(0)
}

func (m *MockWalletUsecase) Display(ctx context.Context, walletID uuid.UUID) (int64, error) {
	args := m.Called(ctx, walletID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockWalletUsecase) CreateWallet(ctx context.Context) (uuid.UUID, error) {
	args := m.Called(ctx)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
