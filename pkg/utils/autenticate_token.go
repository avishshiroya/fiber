package utils

import (
	"fiber/app/models"
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
