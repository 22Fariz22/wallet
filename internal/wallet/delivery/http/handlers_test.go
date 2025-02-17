package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/22Fariz22/wallet/config"
	"github.com/22Fariz22/wallet/internal/wallet"
	"github.com/22Fariz22/wallet/pkg/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDisplayHandler(t *testing.T) {
	e := echo.New()
	mockUsecase := new(wallet.MockWalletUsecase)
	cfg := &config.Config{}
	log := logger.NewMockLogger()
	handler := NewWalletHandler(cfg, mockUsecase, log)

	t.Run("Success", func(t *testing.T) {
		walletID := uuid.New()
		mockUsecase.On("Display", mock.Anything, walletID).Return(int64(1000), nil)

		req := httptest.NewRequest(http.MethodGet, "/wallets/"+walletID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("uuid")
		c.SetParamValues(walletID.String())

		err := handler.Display()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "1000")
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wallets/invalid-uuid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("uuid")
		c.SetParamValues("invalid-uuid")

		err := handler.Display()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid UUID format")
	})

	t.Run("Error fetching balance", func(t *testing.T) {
		walletID := uuid.New()
		mockUsecase.On("Display", mock.Anything, walletID).Return(int64(0), errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/wallets/"+walletID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("uuid")
		c.SetParamValues(walletID.String())

		err := handler.Display()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "internal server error")
		mockUsecase.AssertExpectations(t)
	})
}

func TestOperationHandler(t *testing.T) {
	e := echo.New()
	mockUsecase := new(wallet.MockWalletUsecase)
	cfg := &config.Config{}
	log := logger.NewMockLogger()
	handler := NewWalletHandler(cfg, mockUsecase, log)

	t.Run("Success Deposit", func(t *testing.T) {
		walletID := uuid.New()
		requestBody := map[string]interface{}{
			"walletID":      walletID.String(),
			"operationType": "DEPOSIT",
			"amount":        500,
		}
		jsonData, _ := json.Marshal(requestBody)

		mockUsecase.On("Deposit", mock.Anything, walletID, int64(500)).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Operation()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBufferString(`{"walletID": "not-a-uuid", "operationType": "DEPOSIT", "amount": 100}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Operation()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid UUID format")
	})

	t.Run("Invalid Operation Type", func(t *testing.T) {
		walletID := uuid.New()
		requestBody := map[string]interface{}{
			"walletID":      walletID.String(),
			"operationType": "INVALID_OP",
			"amount":        100,
		}
		jsonData, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Operation()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "operationType must be 'DEPOSIT' or 'WITHDRAW'")
	})

	t.Run("Failed Withdraw", func(t *testing.T) {
		walletID := uuid.New()
		requestBody := map[string]interface{}{
			"walletID":      walletID.String(),
			"operationType": "WITHDRAW",
			"amount":        200,
		}
		jsonData, _ := json.Marshal(requestBody)

		mockUsecase.On("Withdraw", mock.Anything, walletID, int64(200)).Return(errors.New("insufficient funds"))

		req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Operation()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "operation failed")
		mockUsecase.AssertExpectations(t)
	})
}

func TestCreateWalletHandler(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{}
	log := logger.NewMockLogger()

	t.Run("Success", func(t *testing.T) {
		mockUsecase := new(wallet.MockWalletUsecase)
		handler := NewWalletHandler(cfg, mockUsecase, log)

		walletID := uuid.New()
		mockUsecase.On("CreateWallet", mock.Anything).Return(walletID, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/new", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateWallet()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), walletID.String())
		assert.Contains(t, rec.Body.String(), "Wallet created successfully")

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		mockUsecase := new(wallet.MockWalletUsecase)
		handler := NewWalletHandler(cfg, mockUsecase, log)

		mockUsecase.On("CreateWallet", mock.Anything).Return(uuid.UUID{}, errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodPost, "/new", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateWallet()(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "wallet creation failed")

		mockUsecase.AssertExpectations(t)
	})
}
