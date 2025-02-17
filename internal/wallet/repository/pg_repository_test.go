package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestWalletRepo_Display(t *testing.T) {
	db, sqlMock, _ := sqlmock.New()
	defer db.Close()

	mockRedis, redisMock := redismock.NewClientMock()

	logger := logger.NewMockLogger()

	repo := NewWalletRepository(sqlx.NewDb(db, "postgres"), logger, mockRedis)

	t.Run("Success from Redis", func(t *testing.T) {
		walletID := uuid.New()
		cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)

		redisMock.ExpectGet(cacheKey).SetVal("500")

		balance, err := repo.Display(context.Background(), walletID)

		assert.NoError(t, err)
		assert.Equal(t, int64(500), balance)
		redisMock.ExpectationsWereMet()
	})

	t.Run("Success from DB", func(t *testing.T) {
		walletID := uuid.New()
		cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)

		redisMock.ExpectGet(cacheKey).RedisNil()

		sqlMock.ExpectQuery("SELECT amount FROM wallets WHERE wallet_id = ?").
			WithArgs(walletID).WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(500))

		redisMock.ExpectSet(cacheKey, "500", 10*time.Minute).SetVal("OK")
		balance, err := repo.Display(context.Background(), walletID)

		assert.NoError(t, err)
		assert.Equal(t, int64(500), balance)
		redisMock.ExpectationsWereMet()
		sqlMock.ExpectationsWereMet()
	})

	t.Run("Failure from DB", func(t *testing.T) {
		walletID := uuid.New()
		cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)

		redisMock.ExpectGet(cacheKey).RedisNil()
		sqlMock.ExpectQuery(
			"SELECT amount FROM wallets WHERE wallet_id = ?",
		).WithArgs(walletID).WillReturnError(errors.New("db error"))

		balance, err := repo.Display(context.Background(), walletID)

		assert.Error(t, err)
		assert.Equal(t, int64(0), balance)
		sqlMock.ExpectationsWereMet()
	})
}

func TestWalletRepo_Deposit(t *testing.T) {
	db, sqlMock, _ := sqlmock.New()
	defer db.Close()

	mockRedis, redisMock := redismock.NewClientMock()
	logger := logger.NewMockLogger()

	repo := NewWalletRepository(sqlx.NewDb(db, "postgres"), logger, mockRedis)

	t.Run("Success", func(t *testing.T) {
		walletID := uuid.New()
		amount := int64(100)
		newBalance := int64(200)
		cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery("UPDATE wallets SET amount = amount \\+ \\$1 WHERE wallet_id = \\$2 RETURNING amount").
			WithArgs(amount, walletID).
			WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(newBalance))
		sqlMock.ExpectCommit()

		redisMock.ExpectSet(cacheKey, newBalance, 10*time.Minute).SetVal("OK")

		err := repo.Deposit(context.Background(), walletID, amount)

		assert.NoError(t, err)
		assert.Nil(t, sqlMock.ExpectationsWereMet())
		assert.Nil(t, redisMock.ExpectationsWereMet())
	})

	t.Run("DB Failure", func(t *testing.T) {
		walletID := uuid.New()
		amount := int64(100)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery("UPDATE wallets SET amount = amount \\+ \\$1 WHERE wallet_id = \\$2 RETURNING amount").
			WithArgs(amount, walletID).
			WillReturnError(errors.New("db error"))
		sqlMock.ExpectRollback()

		err := repo.Deposit(context.Background(), walletID, amount)

		assert.Error(t, err)
		assert.Equal(t, "failed to update balance: db error", err.Error())
		assert.Nil(t, sqlMock.ExpectationsWereMet())
	})
}

func TestWalletRepo_Withdraw(t *testing.T) {
	db, sqlMock, _ := sqlmock.New()
	defer db.Close()

	mockRedis, redisMock := redismock.NewClientMock()
	logger := logger.NewMockLogger()

	repo := NewWalletRepository(sqlx.NewDb(db, "postgres"), logger, mockRedis)

	t.Run("Success", func(t *testing.T) {
		walletID := uuid.New()
		amount := int64(50)
		newBalance := int64(150)
		cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery("UPDATE wallets SET amount = amount - \\$1 WHERE wallet_id = \\$2 RETURNING amount").
			WithArgs(amount, walletID).
			WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(newBalance))
		sqlMock.ExpectCommit()

		redisMock.ExpectSet(cacheKey, newBalance, 10*time.Minute).SetVal("OK")

		err := repo.Withdraw(context.Background(), walletID, amount)

		assert.NoError(t, err)
		assert.Nil(t, sqlMock.ExpectationsWereMet())
		assert.Nil(t, redisMock.ExpectationsWereMet())
	})

	t.Run("Insufficient Funds", func(t *testing.T) {
		walletID := uuid.New()
		amount := int64(500)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery("UPDATE wallets SET amount = amount - \\$1 WHERE wallet_id = \\$2 RETURNING amount").
			WithArgs(amount, walletID).
			WillReturnError(errors.New("insufficient funds"))
		sqlMock.ExpectRollback()

		err := repo.Withdraw(context.Background(), walletID, amount)

		assert.Error(t, err)
		assert.Equal(t, "insufficient funds", err.Error())
		assert.Nil(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("DB Failure", func(t *testing.T) {
		walletID := uuid.New()
		amount := int64(100)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery("UPDATE wallets SET amount = amount - \\$1 WHERE wallet_id = \\$2 RETURNING amount").
			WithArgs(amount, walletID).
			WillReturnError(errors.New("db error"))
		sqlMock.ExpectRollback()

		err := repo.Withdraw(context.Background(), walletID, amount)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Nil(t, sqlMock.ExpectationsWereMet())
	})
}

func TestWalletRepo_CreateWallet(t *testing.T) {
	db, sqlMock, _ := sqlmock.New()
	defer db.Close()

	mockRedis, redisMock := redismock.NewClientMock()
	logger := logger.NewMockLogger()

	repo := NewWalletRepository(sqlx.NewDb(db, "postgres"), logger, mockRedis)

	t.Run("Success", func(t *testing.T) {
		walletID := uuid.New()
		cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)

		sqlMock.ExpectQuery("INSERT INTO wallets \\(wallet_id, amount\\) VALUES \\(\\$1, 0\\) RETURNING wallet_id").
			WithArgs(walletID).
			WillReturnRows(sqlmock.NewRows([]string{"wallet_id"}).AddRow(walletID))

		redisMock.ExpectSet(cacheKey, 0, 10*time.Minute).SetVal("OK")

		id, err := repo.CreateWallet(context.Background(), walletID)

		assert.NoError(t, err)
		assert.Equal(t, walletID, id)
		assert.Nil(t, sqlMock.ExpectationsWereMet())
		assert.Nil(t, redisMock.ExpectationsWereMet())
	})

	t.Run("DB Failure", func(t *testing.T) {
		walletID := uuid.New()

		sqlMock.ExpectQuery("INSERT INTO wallets \\(wallet_id, amount\\) VALUES \\(\\$1, 0\\) RETURNING wallet_id").
			WithArgs(walletID).
			WillReturnError(errors.New("db error"))

		id, err := repo.CreateWallet(context.Background(), walletID)

		assert.Error(t, err)
		assert.Equal(t, uuid.UUID{}, id)
		assert.Nil(t, sqlMock.ExpectationsWereMet())
	})

	t.Run("Redis Failure", func(t *testing.T) {
		walletID := uuid.New()
		cacheKey := fmt.Sprintf("wallet_balance:%s", walletID)

		sqlMock.ExpectQuery("INSERT INTO wallets \\(wallet_id, amount\\) VALUES \\(\\$1, 0\\) RETURNING wallet_id").
			WithArgs(walletID).
			WillReturnRows(sqlmock.NewRows([]string{"wallet_id"}).AddRow(walletID))

		redisMock.ExpectSet(cacheKey, 0, 10*time.Minute).SetErr(errors.New("redis error"))

		id, err := repo.CreateWallet(context.Background(), walletID)

		assert.NoError(t, err)
		assert.Equal(t, walletID, id)
		assert.Nil(t, sqlMock.ExpectationsWereMet())
		assert.Nil(t, redisMock.ExpectationsWereMet())
	})
}
