package models

import "github.com/google/uuid"

type Wallet struct {
	WalletID uuid.UUID
	Amount   int64
}

