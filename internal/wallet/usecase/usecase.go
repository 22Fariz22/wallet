package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/22Fariz22/wallet/config"
	"github.com/22Fariz22/wallet/internal/wallet"
	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type walletUseCase struct {
	cfg         *config.Config
	walletRepo  wallet.Repository
	redisClient *redis.Client
	logger      logger.Logger
	httpClient  *http.Client
}

func NewWalletUseCase(
	cfg *config.Config,
	walletRepo wallet.Repository,
	redisClient *redis.Client,
	logger logger.Logger) wallet.Usecase {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &walletUseCase{
		cfg:         cfg,
		walletRepo:  walletRepo,
		redisClient: redisClient,
		logger:      logger,
		httpClient:  httpClient,
	}
}

func (u *walletUseCase) Display(ctx context.Context, walletID uuid.UUID) (int64, error) {
	u.logger.Info("Display usecase called")
	return u.walletRepo.Display(ctx, walletID)
}

func (u *walletUseCase) Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error {
	u.logger.Info("Deposit usecase called")
	return u.walletRepo.Deposit(ctx, walletID, amount)
}

func (u *walletUseCase) Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error {
	u.logger.Info("Withdraw usecase called")
	return u.walletRepo.Withdraw(ctx, walletID, amount)
}

func (u *walletUseCase) CreateWallet(ctx context.Context) (uuid.UUID, error) {
	u.logger.Info("CreateWallet usecase called")
	return u.walletRepo.CreateWallet(ctx)
}
