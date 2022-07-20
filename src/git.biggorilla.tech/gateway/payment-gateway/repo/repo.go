package repo

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	sql "database/sql"
	"encoding/base64"
	"fmt"
	"git.biggorilla.tech/gateway/payment-gateway/model"
	"git.biggorilla.tech/gateway/payment-gateway/services/web3"
	_ "github.com/lib/pq"
	"strings"
)

type MerchantRepo interface {
	CreateMerchant(ctx context.Context, name string, email string, user_id string) (interface{}, error)
	GenerateLink(ctx context.Context, id string) (interface{}, error)
	GenerateDepositAddress(ctx context.Context, s services.EthereumService, network string, coin string, user_id string) string
}

type merchantRepo struct {
	db  *sql.DB
	ctx context.Context
}

func NewMerchantRepo(ctx context.Context, db *sql.DB) MerchantRepo {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	return &merchantRepo{
		db,
		ctx,
	}
}

func (r *merchantRepo) GenerateDepositAddress(ctx context.Context, s services.EthereumService, network string, coin string, user_id string) string {
	selectStatment := `SELECT address FROM accounts WHERE user_id='` + user_id + `' AND network ='` + network + `'`
	data1, err1 := r.db.Query(selectStatment)
	if err1 != nil {
		fmt.Println("dsadsa", err1)
	}
	var taddress string
	data1.Next()
	data1.Scan(&taddress)
	fmt.Println(taddress)
	if taddress == "" {
		sqlStatement := `INSERT INTO accounts (user_id, address, private_key, coin, network)
	VALUES ($1, $2, $3, $4, $5) RETURNING address`
		data := s.GenerateNewAddress()
		addresst := ""
		err := r.db.QueryRow(sqlStatement, user_id, data.PublicKey, data.PrivateKey, coin, network).Scan(&addresst)
		if err != nil {
			panic(err)
		}
		return addresst
	}

	return taddress
}

func (r *merchantRepo) GenerateLink(ctx context.Context, id string) (interface{}, error) {
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
		if strings.Contains(err.Error(), "duplicate key") {
			return &model.GenericResponse{
				Code:    409,
				Message: "Already Generated Link",
			}, nil
		}
		return &model.GenericResponse{
			Code:    500,
			Message: "Error Occured",
		}, nil
	}
	return &model.Link{
		ID:       int64(idt),
		PluginID: plugin_idt,
		UserID:   user_idt,
	}, nil
}

func (r *merchantRepo) CreateMerchant(ctx context.Context, name string, email string, user_id string) (interface{}, error) {

	sqlStatement := `INSERT INTO merchants (name, email, user_id) 
	VALUES ($1, $2, $3) RETURNING id, name, email, user_id`
	idt := 0
	emailt := ""
	namet := ""
	user_idt := ""
	err := r.db.QueryRow(sqlStatement, name, email, user_id).Scan(&idt, &namet, &emailt, &user_idt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return &model.GenericResponse{
				Code:    409,
				Message: "Already Created Merchant",
			}, err
		}
		return &model.GenericResponse{
			Code:    500,
			Message: "Error Occured",
		}, err
	}
	return &model.Merchant{
		Name:   namet,
		ID:     int64(idt),
		Email:  emailt,
		UserID: user_id,
	}, nil
}
