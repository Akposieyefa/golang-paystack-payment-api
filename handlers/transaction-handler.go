package handlers

import (
	"akposieyefa/paystack-payment-api/handlers/auth"
	"akposieyefa/paystack-payment-api/helpers"
	"akposieyefa/paystack-payment-api/models"
	"akposieyefa/paystack-payment-api/pkg"
	"akposieyefa/paystack-payment-api/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid request body",
			"success": false,
		})
		return
	}

	reference := helpers.GenerateAccountNumber(12) // Transaction reference

	authUser := auth.AuthUser(r)
	userID := authUser["user"].(models.User).ID
	userEmail := authUser["user"].(models.User).Email

	var wallet models.Wallet
	if err := pkg.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transaction.WalletID = wallet.ID
	transaction.Reference = reference
	transaction.Status = bool(false)

	if err := helpers.Validate(transaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	// Call Paystack to generate payment link
	resp := services.InitializeTransaction(userEmail, int(transaction.Amount), reference)

	url, ok := resp["data"].(map[string]interface{})["authorization_url"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Failed to get payment URL from Paystack response",
			"success": false,
		})
		return
	}

	if err := pkg.DB.Create(&transaction).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Failed to create wallet",
			"success": false,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Payment link generated successfully",
		"data": map[string]interface{}{
			"url": url,
		},
		"success": true,
	})
}

// verify transaction
func VerifyTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var transaction models.Transaction
	var wallet models.Wallet

	paystackCall := services.VerifyTransaction(params["reference"])

	status := paystackCall["data"].(map[string]interface{})["status"]

	if status != true {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Sorry unable to verfy transaction",
			"success": false,
		})
		return
	} else {

		if err := pkg.DB.Where("reference = ?", params["reference"]).First(&transaction).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//credit wallet
		pkg.DB.First(&wallet, transaction.WalletID)

		transaction.Status = true
		pkg.DB.Save(&transaction)

		wallet.Balance = wallet.Balance + transaction.Amount
		pkg.DB.Save(&wallet)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Transaction verified successfully",
			"data":    transaction,
			"success": true,
		})
	}
}
