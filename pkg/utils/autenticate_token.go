package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateToken(uuid uuid.UUID) (string, error) {
	var secretKey = []byte("83f138c1-801b-4f27-bcd6-ee0dca60d349")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  uuid,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		},
	)
	tokenstring, err := token.SignedString(secretKey)
	fmt.Println(tokenstring, err)
	if err != nil {
		return "", err
	}
	return tokenstring, nil
}
