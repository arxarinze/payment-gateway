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
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.Parse(
		accessToken,
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				fmt.Println("error occured 1")
				return nil, fmt.Errorf("unexpected token signing method")
			}
			fmt.Println(manager.secretKey)
			return []byte(manager.secretKey), nil
		},
	)
	fmt.Println(token)
	if err != nil {
		fmt.Println("error occured 2", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		fmt.Println("error occured 3")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
