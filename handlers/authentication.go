package handlers

import (
	"akposieyefa/paystack-payment-api/handlers/auth"
	"akposieyefa/paystack-payment-api/models"
	"encoding/json"
	"net/http"
)

// login user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Invalid request",
			"body":    r.Body,
			"success": false,
		})
		return
	}
	resp := auth.AuthenticateUserDetails(user.Email, user.Password)
	json.NewEncoder(w).Encode(resp)
}

// get authenticated user profile
func UserProfile(w http.ResponseWriter, r *http.Request) {
	resp := auth.AuthUser(r)
	json.NewEncoder(w).Encode(resp)
}
