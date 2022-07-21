package repo

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	sql "database/sql"
	"encoding/base64"
	"fmt"
	"strings"

	"git.biggorilla.tech/gateway/payment-gateway/model"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"git.biggorilla.tech/gateway/payment-gateway/services/web3"
	_ "github.com/lib/pq"
)

type MerchantRepo interface {
	CreateMerchant(ctx context.Context, name string, email string, user_id string) (interface{}, error)
	GenerateLink(ctx context.Context, merchant_id string, user_id string) (interface{}, error)
	GenerateDepositAddress(ctx context.Context, s services.EthereumService, network string, coin string, user_id string) string
	GetPluginLink(ctx context.Context, user_id string, merchant_id string) string
	GetPublicMerchantInfo(ctx context.Context, plugin_id string) (*pb.MerchantPublicResponse, error)
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

func (r *merchantRepo) GetPublicMerchantInfo(ctx context.Context, plugin_id string) (*pb.MerchantPublicResponse, error) {
	selectStatment := `SELECT merchant_id FROM link WHERE plugin_id='` + plugin_id + `'`
	data, err := r.db.Query(selectStatment)
	if err != nil {
		return nil, err
	}
	var merchant_id string
	data.Next()
	data.Scan(&merchant_id)
	selectStatment1 := `SELECT name,email,user_id FROM merchants WHERE id='` + merchant_id + `'`
	data1, err1 := r.db.Query(selectStatment1)
	if err1 != nil {
		return nil, err
	}
	var email string
	var name string
	var user_id string
	data1.Next()
	data1.Scan(&name, &email, &user_id)
	return &pb.MerchantPublicResponse{
		Name:       name,
		Email:      email,
		MerchantId: merchant_id,
		UserId:     user_id,
	}, nil
}

func (r *merchantRepo) GetPluginLink(ctx context.Context, user_id string, merchant_id string) string {
	selectStatment := `SELECT plugin_id FROM link WHERE user_id='` + user_id + `' AND merchant_id ='` + merchant_id + `'`
	data, err := r.db.Query(selectStatment)
	if err != nil {
		return err.Error()
	}
	var plugin_id string
	data.Next()
	data.Scan(&plugin_id)

	return "$BASE_HOST" + "/donate/" + plugin_id
}
func (r *merchantRepo) GenerateDepositAddress(ctx context.Context, s services.EthereumService, network string, coin string, user_id string) string {
	selectStatment := `SELECT address FROM accounts WHERE user_id='` + user_id + `' AND network ='` + network + `'`
	data1, err1 := r.db.Query(selectStatment)
	if err1 != nil {
		return err1.Error()
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

func (r *merchantRepo) GenerateLink(ctx context.Context, merchant_id string, user_id string) (interface{}, error) {
	data := []byte(merchant_id)
	hash := sha256.Sum256(data)
	plugin_id := fmt.Sprint(hash)
	md5hash := md5.Sum([]byte(plugin_id))
	plugin_id = base64.URLEncoding.EncodeToString(md5hash[:])
	sqlStatement := `INSERT INTO link (plugin_id, user_id, merchant_id) 
	VALUES ($1, $2, $3) RETURNING id, plugin_id, user_id, merchant_id`
	idt := 0
	plugin_idt := ""
	user_idt := ""
	merchant_idt := ""
	err := r.db.QueryRow(sqlStatement, plugin_id, user_id, merchant_id).Scan(&idt, &plugin_idt, &user_idt, &merchant_idt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return &model.GenericResponse{
				Code:    409,
				Message: "Already Generated Link",
			}, err
		}
		return &model.GenericResponse{
			Code:    500,
			Message: "Error Occured" + err.Error(),
		}, err
	}
	return &model.Link{
		ID:         int64(idt),
		PluginID:   plugin_idt,
		UserID:     user_idt,
		MerchantID: merchant_idt,
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
