package helpers

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

func (manager *JWTManager) Verify(accessToken string) (*jwt.StandardClaims, error) {
	fmt.Println(manager.secretKey, accessToken)
	token, err := jwt.ParseWithClaims(
		accessToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				fmt.Println("error occured 1")
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		fmt.Println("error occured 2")
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	fmt.Println("this is token", token)
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		fmt.Println("error occured 3")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
