package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Balance      uint   `json:"balance" validate:"required"`
	AcctNumber   string `json:"account_number" validate:"required"`
	UserID       uint   `json:"userId"`
	Transactions []Transaction
}
