package helpers

import (
	"context"
	"fmt"
	_ "fmt"
)

type Identity interface {
	GetIdentity( /*auth metadata.MD*/ ctx context.Context) string
}

type identity struct {
}

func NewIdentity(ctx context.Context) Identity {
	return &identity{}
}

func (i *identity) GetIdentity(ctx context.Context) string {
	identity := ctx.Value("user").(map[string]interface{})["data"].(map[string]interface{})["id"]
	id := fmt.Sprintf("%v", identity)
	// a := auth.Get("authorization")[0]
	// token := strings.Split(a, " ")[1]
	// tk, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
	// 	return t.Valid, nil
	// })
	// data := tk.Claims.(jwt.MapClaims)
	// username := fmt.Sprint(data["username"])
	return id
}
