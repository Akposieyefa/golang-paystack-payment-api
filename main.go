package main

import (
	"akposieyefa/paystack-payment-api/pkg"
	"akposieyefa/paystack-payment-api/routers"
)

func main() {
	pkg.ConnectToDB()
	routers.Router()
}
