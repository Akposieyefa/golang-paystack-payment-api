package handlers

import (
	"akposieyefa/paystack-payment-api/handlers/auth"
	"akposieyefa/paystack-payment-api/helpers"
	"akposieyefa/paystack-payment-api/models"
	"akposieyefa/paystack-payment-api/pkg"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	json.NewDecoder(r.Body).Decode(user)

	err := helpers.Validate(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Password Encryption  failed",
			"success": false,
		})
		return
	}

	user.Password = string(pass)

	createdUser := pkg.DB.Create(user)
	var errMessage = createdUser.Error

	if createdUser.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": errMessage,
			"success": false,
		})
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"users":   user,
		"success": true,
	})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user models.User
	pkg.DB.First(&user, params["id"])
	json.NewDecoder(r.Body).Decode(&user)

	authUser := auth.AuthUser(r)
	authID := authUser["user"].(models.User).ID

	if user.ID != authID {
		w.WriteHeader(http.StatusNotAcceptable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Sorry you are not allowed to update this information",
			"success": false,
		})
		return
	}

	pkg.DB.Save(&user)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User details updated successfully",
		"data":    user,
		"success": true,
	})
}
