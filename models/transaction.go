package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Title       string `json:"title" validate:"required"`
	Amount      uint   `json:"amount" validate:"required"`
	Status      bool   `json:"status"`
	Description string `json:"description" validate:"required"`
	Reference   string `json:"reference" validate:"required"`
	WalletID    uint   `json:"walletId"`
}
