package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthDto struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Jti      uuid.UUID `json:"jti"`
}

func CreateToken(auth AuthDto) (string, error) {
	var secretKey = []byte("83f138c1-801b-4f27-bcd6-ee0dca60d349")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       auth.ID,
			"email":    auth.Email,
			"username": auth.Username,
			"jti":      auth.Jti,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		},
	)
	tokenstring, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}
func CreateRefreshToken(auth AuthDto) (string, error) {
	var secretKey = []byte("0b60783b-783f-472a-a5ba-b9759d47e0c0")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       auth.ID,
			"email":    auth.Email,
			"username": auth.Username,
			"jti":      auth.Jti,
			"exp":      time.Now().Add(time.Hour * 24 * 6).Unix(),
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
func VerifyRefreshToken(tokenString string) (jwt.MapClaims, error) {
	var secretKey = []byte("0b60783b-783f-472a-a5ba-b9759d47e0c0")

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
