package services

import (
	"akposieyefa/paystack-payment-api/helpers"
	"akposieyefa/paystack-payment-api/pkg"
	"bytes"
	"encoding/json"
	"net/http"
)

var apiKey = pkg.LoadEnv("PAYSTACK_SECRET_KEY")
var apiUrl = pkg.LoadEnv("PAYSTACK_URL")
var callBackUrl = pkg.LoadEnv("PAYSTACK_CALLBACK_URL")

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

	req, err := http.NewRequest("POST", apiUrl+"/transaction/initialize", bytes.NewBuffer(jsonData))
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

	req, err := http.NewRequest("GET", apiUrl+"/transaction/verify/"+reference, nil)
	if err != nil {
		return map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		}
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
