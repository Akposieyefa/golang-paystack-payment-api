package services

import (
	"akposieyefa/paystack-payment-api/helpers"
	"akposieyefa/paystack-payment-api/pkg"
	"bytes"
	"encoding/json"
	"net/http"
)

var apiKey = pkg.LoadEnv("PAYSTACK_SECRET_KEY")
var callBackUrl = pkg.LoadEnv("PAYSTACK_CALLBACK_UR")

// initialize paystack payment service
func InitializeTransaction(email string, amount int, reference string) map[string]interface{} {

	params := map[string]interface{}{
		"email":        email,
		"amount":       helpers.ConvertToKobo(amount),
		"reference":    reference,
		"callback_url": callBackUrl + "=" + reference,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		}
	}

	url := "https://api.paystack.co/transaction/initialize"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		}
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		}
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		}
	}

	return map[string]interface{}{
		"data": data["data"],
	}
}

// verify transaction service
func VerifyTransaction(reference string) map[string]interface{} {

	url := "https://api.paystack.co/transaction/verify/" + reference
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		}
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		}
	}

	return map[string]interface{}{
		"data": data,
	}
}
