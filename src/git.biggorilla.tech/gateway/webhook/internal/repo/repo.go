package repo

import (
	"context"
	sql "database/sql"
	"fmt"

	"git.biggorilla.tech/gateway/webhook/internal/models"
	_ "github.com/lib/pq"
)

type WebhookRepo interface {
	CheckForAddress(ctx context.Context, address string) (bool, error)
	InsertTransaction(ctx context.Context, data models.Transaction) (bool, error)
	GetCoinForNetwork(ctx context.Context, address string) (*models.Asset, error)
}

type webhookRepo struct {
	db  *sql.DB
	ctx context.Context
}

// GetCoinForNetwork implements WebhookRepo
func (r *webhookRepo) GetCoinForNetwork(ctx context.Context, address string) (*models.Asset, error) {
	selectStatment := `SELECT coin, network FROM assets WHERE address='` + address + `'`
	data, err := r.db.Query(selectStatment)
	if err != nil {
		fmt.Print(err)
		return &models.Asset{}, err
	}
	defer data.Close()
	var coin string
	var network string
	data.Next()
	data.Scan(&coin, &network)
	if coin != "" {
		return &models.Asset{}, nil
	}
	return &models.Asset{
		Coin:    coin,
		Network: network,
	}, nil
}

// CheckForAddress implements WebhookRepo
func (r *webhookRepo) CheckForAddress(ctx context.Context, address string) (bool, error) {
	selectStatment := `SELECT user_id, merchant_id FROM accounts WHERE address='` + address + `'`
	data, err := r.db.Query(selectStatment)
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	defer data.Close()
	var user_id string
	var merchant_id string
	data.Next()
	data.Scan(&user_id, &merchant_id)
	if user_id != "" {
		return true, nil
	}
	return false, nil
}

// InsertTransaction implements WebhookRepo
func (r *webhookRepo) InsertTransaction(ctx context.Context, data models.Transaction) (bool, error) {
	selectStatment1 := `SELECT user_id, merchant_id FROM accounts WHERE address='` + data.To + `'`
	data1, err := r.db.Query(selectStatment1)
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	defer data1.Close()
	var user_id string
	var merchant_id string
	data1.Next()
	data1.Scan(&user_id, &merchant_id)
	fmt.Println(user_id, merchant_id, data.From, data.To, data.TxHash, data.Value)
	sqlStatement := `INSERT INTO transactions (tx_hash,sender,reciever,value,type,status,merchant_id, user_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	dat, err := r.db.Exec(sqlStatement, data.TxHash, data.From, data.To, data.Value, "credit", true, merchant_id, user_id)
	fmt.Println(dat)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, nil
}

func NewWebhookRepo(ctx context.Context, db *sql.DB) WebhookRepo {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	return &webhookRepo{
		db,
		ctx,
	}
}
