package repo

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	sql "database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"git.biggorilla.tech/gateway/payment-gateway/model"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	services "git.biggorilla.tech/gateway/payment-gateway/services/web3"
	_ "github.com/lib/pq"
)

type MerchantRepo interface {
	CreateMerchant(ctx context.Context, name string, email string, address string, avatar string, user_id string) (interface{}, error)
	UpdateMerchant(ctx context.Context, name string, email string, address string, avatar string, user_id string, merchant_id int64) (interface{}, error)
	GetMerchants(ctx context.Context, user_id string) (*[]model.Merchant, error)
	GenerateLink(ctx context.Context, merchant_id int64, user_id string) (interface{}, error)
	GenerateDepositAddress(ctx context.Context, s services.EthereumService, network string, coin string, plugin_id string) (string, error)
	GetPluginLink(ctx context.Context, user_id string, merchant_id string, typeOf string) (string, error)
	GetPublicMerchantInfo(ctx context.Context, plugin_id string) (*pb.MerchantPublicResponse, error)
	GetTransactions(ctx context.Context, merchant_id string, user_id string) (*[]model.Transaction, error)
}

type merchantRepo struct {
	db  *sql.DB
	ctx context.Context
}

// GetTransactions implements MerchantRepo
func (m *merchantRepo) GetTransactions(ctx context.Context, merchant_id string, user_id string) (*[]model.Transaction, error) {
	result := []model.Transaction{}
	selectStatment := `SELECT tx_hash, sender, reciever, value FROM transactions WHERE user_id='` + user_id + `' AND merchant_id='` + merchant_id + `'`
	data, err := m.db.Query(selectStatment)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	for data.Next() {
		dataSet := model.Transaction{
			TxHash: "",
			From:   "",
			To:     "",
			Value:  "",
		}
		err = data.Scan(&dataSet.TxHash, &dataSet.From, &dataSet.To, &dataSet.Value)
		if err != nil {
			// handle this error
			panic(err)
		}
		result = append(result, dataSet)
	}
	return &result, nil
}

// UpdateMerchant implements MerchantRepo
func (m *merchantRepo) UpdateMerchant(ctx context.Context, name string, email string, address string, avatar string, user_id string, merchant_id int64) (interface{}, error) {
	sqlStatement := `UPDATE merchants SET name = $2, email = $3, address=$4, avatar=$5 WHERE user_id = $1 AND id=$6;`
	res, err := m.db.Exec(sqlStatement, user_id, name, email, address, avatar, merchant_id)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	return count, nil
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

// GetMerchants implements MerchantRepo
func (m *merchantRepo) GetMerchants(ctx context.Context, user_id string) (*[]model.Merchant, error) {
	result := []model.Merchant{}
	selectStatment := `SELECT id, name, email, user_id, avatar, address FROM merchants WHERE user_id='` + user_id + `'`
	data, err := m.db.Query(selectStatment)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	for data.Next() {
		dataSet := model.Merchant{
			ID:      0,
			Name:    "",
			Email:   "",
			UserID:  user_id,
			Avatar:  "",
			Address: "",
		}
		err = data.Scan(&dataSet.ID, &dataSet.Name, &dataSet.Email, &dataSet.UserID, &dataSet.Avatar, &dataSet.Address)
		if err != nil {
			// handle this error
			panic(err)
		}
		result = append(result, dataSet)
	}
	return &result, nil
}
func (r *merchantRepo) GetPublicMerchantInfo(ctx context.Context, plugin_id string) (*pb.MerchantPublicResponse, error) {
	selectStatment := `SELECT merchant_id FROM link WHERE plugin_id='` + plugin_id + `'`
	data, err := r.db.Query(selectStatment)
	if err != nil {
		return nil, err
	}
	defer data.Close()
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

func (r *merchantRepo) GetPluginLink(ctx context.Context, user_id string, merchant_id string, typeOf string) (string, error) {
	selectStatment := `SELECT plugin_id FROM link WHERE user_id='` + user_id + `' AND merchant_id ='` + merchant_id + `'`
	data, err := r.db.Query(selectStatment)
	if err != nil {
		return "", err
	}
	defer data.Close()
	var plugin_id string
	data.Next()
	data.Scan(&plugin_id)
	if plugin_id == "" {
		return "", errors.New("error")
	}
	if strings.ToLower(typeOf) == "iframe" {
		re := regexp.MustCompile(`\t?\r?\n`)
		input := `<iframe src='http://localhost:3000/payment-gateway/` + plugin_id + `' style='height: 600px;width: 300px;'></iframe>`
		input = re.ReplaceAllString(input, "")
		return input, nil
	}
	return "$BASE_HOST" + "/payment-gateway/" + plugin_id, nil
}
func (r *merchantRepo) GenerateDepositAddress(ctx context.Context, s services.EthereumService, network string, coin string, plugin_id string) (string, error) {
	if plugin_id == "" {
		return "", errors.New("missing property plugin_id")
	}
	selectStatment := `SELECT user_id, merchant_id FROM link WHERE plugin_id='` + plugin_id + `'`
	data, err := r.db.Query(selectStatment)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	defer data.Close()
	var user_id string
	var merchant_id string
	data.Next()
	data.Scan(&user_id, &merchant_id)
	fmt.Println(merchant_id, user_id)
	if user_id == "" {
		return "", errors.New("plugin_id is invalid")
	}
	selectStatment1 := `SELECT address FROM accounts WHERE user_id='` + user_id + `' AND network ='` + network + `' AND merchant_id='` + merchant_id + `'`
	data1, err1 := r.db.Query(selectStatment1)
	if err1 != nil {
		return "", err1
	}
	var taddress string
	data1.Next()
	data1.Scan(&taddress)
	if taddress == "" {
		sqlStatement := `INSERT INTO accounts (user_id,merchant_id, address, private_key, coin, network)
	VALUES ($1, $2, $3, $4, $5,$6) RETURNING address`
		data := s.GenerateNewAddress()
		addresst := ""
		err := r.db.QueryRow(sqlStatement, user_id, merchant_id, data.PublicKey, data.PrivateKey, coin, network).Scan(&addresst)
		if err != nil {
			panic(err)
		}
		return addresst, nil
	}

	return taddress, nil
}

func (r *merchantRepo) GenerateLink(ctx context.Context, merchant_id int64, user_id string) (interface{}, error) {
	fmt.Println(merchant_id, user_id)
	selectStatment := `SELECT id FROM merchants WHERE user_id='` + user_id + `' AND id=` + strconv.Itoa(int(merchant_id))
	dta, err0 := r.db.Query(selectStatment)
	if err0 != nil {
		fmt.Print(err0)
		return "", err0
	}
	fmt.Println(dta)
	defer dta.Close()
	var idtemp string
	dta.Next()
	dta.Scan(&idtemp)
	if idtemp == "" {
		return nil, errors.New("Forbidden")
	}
	data := []byte(strconv.Itoa(int(merchant_id)) + user_id + time.Now().GoString())
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
		fmt.Println(err)
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

func (r *merchantRepo) CreateMerchant(ctx context.Context, name string, email string, address string, avatar string, user_id string) (interface{}, error) {

	sqlStatement := `INSERT INTO merchants (name, email, user_id, address, avatar) 
	VALUES ($1, $2, $3, $4,$5) RETURNING id, name, email, user_id, address, avatar`
	idt := 0
	emailt := ""
	namet := ""
	user_idt := ""
	addresst := ""
	avatart := ""
	err := r.db.QueryRow(sqlStatement, name, email, user_id, address, avatar).Scan(&idt, &namet, &emailt, &user_idt, &addresst, &avatart)
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
		Name:    namet,
		ID:      int64(idt),
		Email:   emailt,
		UserID:  user_id,
		Avatar:  avatart,
		Address: addresst,
	}, nil
}
