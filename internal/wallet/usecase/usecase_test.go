package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/22Fariz22/wallet/config"
	"github.com/22Fariz22/wallet/internal/wallet/usecase"
	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок репозитория
type MockWalletRepo struct {
	mock.Mock
}

func (m *MockWalletRepo) Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error {
	args := m.Called(ctx, walletID, amount)
	return args.Error(0)
}

func (m *MockWalletRepo) Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error {
	args := m.Called(ctx, walletID, amount)
	return args.Error(0)
}

func (m *MockWalletRepo) Display(ctx context.Context, walletID uuid.UUID) (int64, error) {
	args := m.Called(ctx, walletID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockWalletRepo) CreateWallet(ctx context.Context, walletID uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func TestWalletUseCase(t *testing.T) {
	mockRepo := new(MockWalletRepo)
	mockLogger := logger.NewMockLogger()
	cfg := &config.Config{}
	useCase := usecase.NewWalletUseCase(cfg, mockRepo, nil, mockLogger)

	ctx := context.Background()
	walletID := uuid.New()
	amount := int64(100)

	t.Run("Deposit Success", func(t *testing.T) {
		mockRepo.On("Deposit", ctx, walletID, amount).Return(nil).Once()

		err := useCase.Deposit(ctx, walletID, amount)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deposit Error", func(t *testing.T) {
		mockRepo.On("Deposit", ctx, walletID, amount).Return(errors.New("deposit error")).Once()

		err := useCase.Deposit(ctx, walletID, amount)

		assert.EqualError(t, err, "deposit error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Withdraw Success", func(t *testing.T) {
		mockRepo.On("Withdraw", ctx, walletID, amount).Return(nil).Once()

		err := useCase.Withdraw(ctx, walletID, amount)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Withdraw Error", func(t *testing.T) {
		mockRepo.On("Withdraw", ctx, walletID, amount).Return(errors.New("withdraw error")).Once()

		err := useCase.Withdraw(ctx, walletID, amount)

		assert.EqualError(t, err, "withdraw error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Display Success", func(t *testing.T) {
		mockRepo.On("Display", ctx, walletID).Return(amount, nil).Once()

		balance, err := useCase.Display(ctx, walletID)

		assert.NoError(t, err)
		assert.Equal(t, amount, balance)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Display Error", func(t *testing.T) {
		mockRepo.On("Display", ctx, walletID).Return(int64(0), errors.New("display error")).Once()

		balance, err := useCase.Display(ctx, walletID)

		assert.EqualError(t, err, "display error")
		assert.Equal(t, int64(0), balance)
		mockRepo.AssertExpectations(t)
	})
	t.Run("CreateWallet Success", func(t *testing.T) {
		newWalletID := uuid.New()
		mockRepo.On("CreateWallet", ctx).Return(newWalletID, nil).Once()

		resultID, err := useCase.CreateWallet(ctx)

		assert.NoError(t, err)
		assert.Equal(t, newWalletID, resultID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CreateWallet Error", func(t *testing.T) {
		mockRepo.On("CreateWallet", ctx).Return(uuid.Nil, errors.New("create wallet error")).Once()

		resultID, err := useCase.CreateWallet(ctx)

		assert.EqualError(t, err, "create wallet error")
		assert.Equal(t, uuid.Nil, resultID)
		mockRepo.AssertExpectations(t)
	})
}
