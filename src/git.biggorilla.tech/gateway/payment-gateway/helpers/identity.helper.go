package helpers

import (
	"context"
	"fmt"
	_ "fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

type Identity interface {
	GetIdentity(auth metadata.MD) string
}

type identity struct {
}

func NewIdentity(ctx context.Context) Identity {
	return &identity{}
}

func (i *identity) GetIdentity(auth metadata.MD) string {
	a := auth.Get("authorization")[0]
	token := strings.Split(a, " ")[1]
	tk, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return t.Valid, nil
	})
	data := tk.Claims.(jwt.MapClaims)
	username := fmt.Sprint(data["username"])
	return username
}
