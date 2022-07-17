package repo

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	sql "database/sql"
	"encoding/base64"
	"fmt"

	"git.biggorilla.tech/gateway/payment-gateway/model"
	_ "github.com/lib/pq"
)

type MerchantRepo interface {
	CreateMerchant(ctx context.Context, name string, email string, user_id string) (*model.Merchant, error)
	GenerateLink(ctx context.Context, id string) (*model.Link, error)
}

type merchantRepo struct {
	db *sql.DB
}

func NewMerchantRepo(ctx context.Context, db *sql.DB) MerchantRepo {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	return &merchantRepo{
		db,
	}
}

func (r *merchantRepo) GenerateLink(ctx context.Context, id string) (*model.Link, error) {
	data := []byte(id)
	hash := sha256.Sum256(data)
	plugin_id := fmt.Sprint(hash)
	md5hash := md5.Sum([]byte(plugin_id))
	plugin_id = base64.URLEncoding.EncodeToString(md5hash[:])
	sqlStatement := `INSERT INTO link (plugin_id, user_id) 
	VALUES ($1, $2) RETURNING id, plugin_id, user_id`
	idt := 0
	plugin_idt := ""
	user_idt := ""
	err := r.db.QueryRow(sqlStatement, plugin_id, id).Scan(&idt, &plugin_idt, &user_idt)
	if err != nil {
		panic(err)
	}
	return &model.Link{
		ID:       int64(idt),
		PluginID: plugin_idt,
		UserID:   user_idt,
	}, nil
}

func (r *merchantRepo) CreateMerchant(ctx context.Context, name string, email string, user_id string) (*model.Merchant, error) {

	sqlStatement := `INSERT INTO merchants (name, email, user_id) 
	VALUES ($1, $2, $3) RETURNING id, name, email, user_id`
	idt := 0
	emailt := ""
	namet := ""
	user_idt := ""
	err := r.db.QueryRow(sqlStatement, name, email, user_id).Scan(&idt, &namet, &emailt, &user_idt)
	if err != nil {
		panic(err)
	}
	return &model.Merchant{
		Name:   namet,
		ID:     int64(idt),
		Email:  emailt,
		UserID: user_id,
	}, nil
}
