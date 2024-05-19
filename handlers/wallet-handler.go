package handlers

import (
	"akposieyefa/paystack-payment-api/handlers/auth"
	"akposieyefa/paystack-payment-api/helpers"
	"akposieyefa/paystack-payment-api/models"
	"akposieyefa/paystack-payment-api/pkg"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// create wallet
func CreateWallet(w http.ResponseWriter, r *http.Request) {
	var wallet models.Wallet

	authUser := auth.AuthUser(r)
	authUserID := authUser["user"].(models.User).ID

	if err := pkg.DB.Where("user_id = ?", authUserID).First(&wallet).Error; err == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Sorry, user cannot have more than one wallet",
			"success": false,
		})
		return
	}

	wallet.UserID = authUserID
	wallet.AcctNumber = helpers.GenerateAccountNumber(10)
	wallet.Balance = 1000

	json.NewDecoder(r.Body).Decode(&wallet)

	if err := helpers.Validate(wallet); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	if err := pkg.DB.Create(&wallet).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Failed to create wallet",
			"success": false,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Wallet created successfully",
		"data":    wallet,
		"success": true,
	})
}

// get all user wallets
func GetAllMyWallets(w http.ResponseWriter, r *http.Request) {

	authUser := auth.AuthUser(r)
	userID := authUser["user"].(models.User).ID

	var wallets []models.Wallet
	if err := pkg.DB.Where("user_id = ?", userID).Find(&wallets).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User wallets fetched successfully",
		"data":    wallets,
		"success": true,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// delete wallet
func DeleteWallet(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	var wallet models.Wallet
	pkg.DB.First(&wallet, param["id"])

	authUser := auth.AuthUser(r)
	userID := authUser["user"].(models.User).ID

	if wallet.UserID != userID {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Sorry you can not delete other users wallet",
			"success": false,
		})
		return
	}
	pkg.DB.Delete(&wallet)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Wallet deleted successfully",
		"data":    wallet,
		"success": true,
	})
}
