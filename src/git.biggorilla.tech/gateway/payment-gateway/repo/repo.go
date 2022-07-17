package repo

import (
	"context"
	sql "database/sql"

	"git.biggorilla.tech/gateway/payment-gateway/model"
	_ "github.com/lib/pq"
)

type MerchantRepo interface {
	CreateMerchant(ctx context.Context, name string, email string, user_id string) (*model.Merchant, error)
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
