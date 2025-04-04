package utils

import (
	"fiber/app/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateToken(user models.User) (string, error) {
	var secretKey = []byte("83f138c1-801b-4f27-bcd6-ee0dca60d349")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"jti":      uuid.New(),
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		},
	)
	tokenstring, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	var secretKey = []byte("83f138c1-801b-4f27-bcd6-ee0dca60d349")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ensure the signing method is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
