package utils

import (
	"final-project/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, userRole string) (string, error) {
	fmt.Println("Generating token for user:", userID, "with role:", userRole)
	claims := jwt.MapClaims{
		"user_id": userID,
		"user_role": userRole,
		"exp": time.Now().Add(config.GetJWTExpirationDuration()).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(config.GetJWTSecret())
}

func ValidateToken(tokenString string) (uint, string, error) {
	// parsing token, return token 
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		return config.GetJWTSecret(), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := uint(claims["user_id"].(float64))
		userRole := claims["user_role"].(string)
		return userId, userRole, nil
	}

	return 0, "", err

}

