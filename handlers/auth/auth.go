package auth

import (
	"akposieyefa/paystack-payment-api/models"
	"akposieyefa/paystack-payment-api/pkg"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(pkg.LoadEnv("JWT_SECRET"))

// authenticate user details
func AuthenticateUserDetails(email, password string) map[string]interface{} {
	user := &models.User{}

	if err := pkg.DB.Where("Email = ?", email).First(user).Error; err != nil {
		return map[string]interface{}{
			"message": "Email address not found",
			"status":  false,
		}
	}
	expirationTime := time.Now().Add(time.Minute * 100000).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return map[string]interface{}{
			"message": "Invalid login credentials. Please try again",
			"status":  false,
		}
	}

	tk := &models.Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(os.Getenv("JWT_ALGO")), tk)

	tokenString, error := token.SignedString(jwtKey)
	if error != nil {
		return map[string]interface{}{
			"message": error.Error(),
			"status":  false,
		}
	}

	return map[string]interface{}{
		"message": "User logged in successfully",
		"data": map[string]interface{}{
			"token": tokenString,
			"user":  user,
		},
		"status": true,
	}
}

// get authenticated user
func AuthUser(request *http.Request) map[string]interface{} {
	tokenString := request.Header.Get("Authorization")

	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return map[string]interface{}{
				"message": "Signature validation error",
				"status":  false,
			}
		}
		return map[string]interface{}{
			"message": "Signature validation error",
			"status":  false,
		}
	}

	if !tkn.Valid {
		return map[string]interface{}{
			"message": "Invalid token",
			"status":  false,
		}
	}
	loggedIn, err := getUserByEmail(claims.Email)
	if err != nil {
		return map[string]interface{}{
			"message": "Error retrieving user data",
			"status":  false,
		}
	}
	return map[string]interface{}{
		"message": "User profile pulled successfully",
		"user":    loggedIn,
		"success": true,
	}
}

// get user by email
func getUserByEmail(email string) (models.User, error) {
	var user models.User
	if err := pkg.DB.Where("Email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
