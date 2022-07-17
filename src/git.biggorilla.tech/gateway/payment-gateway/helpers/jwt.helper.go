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

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
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
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		fmt.Println("error occured 3")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
