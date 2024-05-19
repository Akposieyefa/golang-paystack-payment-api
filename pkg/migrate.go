package pkg

import (
	"akposieyefa/paystack-payment-api/models"
)

func MigrateTables() {
	ConnectToDB()
	DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})
}
