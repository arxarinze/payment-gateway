package repo

import (
	"context"
	sql "database/sql"

	"git.biggorilla.tech/gateway/payment-gateway/model"
	_ "github.com/lib/pq"
)

type MerchantRepo interface {
	CreateMerchant(ctx context.Context, name string, email string) (*model.Merchant, error)
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

func (r *merchantRepo) CreateMerchant(ctx context.Context, name string, email string) (*model.Merchant, error) {

	sqlStatement := `INSERT INTO merchants (name, email) 
	VALUES ($1, $2) RETURNING id, name, email`
	idt := 0
	emailt := ""
	namet := ""
	err := r.db.QueryRow(sqlStatement, name, email).Scan(&idt, &namet, &emailt)
	if err != nil {
		panic(err)
	}
	return &model.Merchant{
		Name:  namet,
		ID:    int64(idt),
		Email: emailt,
	}, nil
}
