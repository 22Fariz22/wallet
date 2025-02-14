package models

import "github.com/google/uuid"

type Wallet struct {
	WalletID uuid.UUID `json:"wallet_id" gorm:"type:uuid;primaryKey" db:"wallet_id"`
	Amount   int64     `json:"amount" gorm:"not null" db:"amount"`
}

