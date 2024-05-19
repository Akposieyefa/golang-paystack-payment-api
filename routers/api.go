package routers

import (
	"akposieyefa/paystack-payment-api/handlers"
	"akposieyefa/paystack-payment-api/middlewares"
	"akposieyefa/paystack-payment-api/pkg"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var port = pkg.LoadEnv("APP_PORT")

func Router() {
	router := mux.NewRouter()
	router.Use(middlewares.Middleware)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Welcome to simple user wallet management system using paystack",
			"success": true,
		})
	}).Methods("GET")

	apiRoutes := router.PathPrefix("/api/v1").Subrouter()

	apiRoutes.HandleFunc("/auth/login", handlers.LoginUser).Methods("POST")
	apiRoutes.HandleFunc("/auth/profiles", handlers.UserProfile).Methods("GET")
	apiRoutes.HandleFunc("/users/create", handlers.CreateUser).Methods("POST")
	apiRoutes.HandleFunc("/users/update/{id}", handlers.UpdateUser).Methods("PATCH")
	apiRoutes.HandleFunc("/wallets/create", handlers.CreateWallet).Methods("POST")
	apiRoutes.HandleFunc("/wallets/all", handlers.GetAllMyWallets).Methods("GET")
	apiRoutes.HandleFunc("/wallets/delete/{id}", handlers.DeleteWallet).Methods("DELETE")
	apiRoutes.HandleFunc("/transactions/create", handlers.CreateTransaction).Methods("POST")
	apiRoutes.HandleFunc("/transactions/verify/{reference}", handlers.VerifyTransaction).Methods("GET")

	s := &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    120 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
