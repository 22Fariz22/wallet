package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/22Fariz22/wallet/internal/wallet"
	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type walletRepo struct {
	db          *sqlx.DB
	logger      logger.Logger
	redisClient *redis.Client
}

func NewWalletRepository(db *sqlx.DB, logger logger.Logger, redisClient *redis.Client) wallet.Repository {
	return &walletRepo{db: db, logger: logger, redisClient: redisClient}
}

func (r *walletRepo) Display(ctx context.Context, walletID uuid.UUID) (int64, error) {
	r.logger.Info("Display repo called")
	cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)

	//  Пытаемся получить баланс из Redis
	balanceStr, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		balance, _ := strconv.ParseInt(balanceStr, 10, 64)
		return balance, nil
	}

	//  Если в Redis нет, идем в БД
	var balance int64
	err = r.db.GetContext(ctx, &balance, "SELECT amount FROM wallets WHERE wallet_id = $1", walletID)
	if err != nil {
		return 0, err
	}

	//  Обновляем кэш в Redis
	r.redisClient.Set(ctx, cacheKey, balance, 10*time.Minute)

	return balance, nil
}

func (r *walletRepo) Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error {
	r.logger.Info("Display repo called")
	var newBalance int64

	// 1. Обновляем баланс в БД и сразу получаем новое значение
	err := r.db.GetContext(ctx, &newBalance,
		"UPDATE wallets SET amount = amount + $1 WHERE wallet_id = $2 RETURNING amount",
		amount, walletID)
	if err != nil {
		return err
	}

	// 2. Обновляем кэш в Redis
	cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)
	r.redisClient.Set(ctx, cacheKey, newBalance, 10*time.Minute)

	r.logger.Infof("Deposit success: wallet %s, new balance: %d", walletID, newBalance)
	return nil
}

func (r *walletRepo) Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error {
	var newBalance int64

	// Обновляем баланс, но не даем уйти в минус
	err := r.db.GetContext(ctx, &newBalance,
		"UPDATE wallets SET amount = amount - $1 WHERE wallet_id = $2 AND amount >= $1 RETURNING amount",
		amount, walletID)
	if err != nil {
		// Если ошибка из-за нехватки средств (нет RETURNING), возвращаем ошибку
		if err == sql.ErrNoRows {
			return fmt.Errorf("insufficient funds")
		}
		return err
	}

	// Обновляем кэш в Redis
	cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)
	r.redisClient.Set(ctx, cacheKey, newBalance, 10*time.Minute)

	r.logger.Infof("Withdraw success: wallet %s, new balance: %d", walletID, newBalance)
	return nil
}

// CreateWallet создаем новый кошелек с нулевым балансом
func (r *walletRepo) CreateWallet(ctx context.Context) (uuid.UUID, error) {
	r.logger.Info("CreateWallet repo called")

	walletID := uuid.New()

	query := `INSERT INTO wallets (wallet_id, amount) VALUES ($1, 0) RETURNING wallet_id`
	err := r.db.QueryRowContext(ctx, query, walletID).Scan(&walletID)
	if err != nil {
		r.logger.Errorf("Failed to create wallet: %v", err)
		return uuid.UUID{}, err
	}

	r.logger.Infof("Wallet created successfully: %s", walletID)
	return walletID, nil
}
