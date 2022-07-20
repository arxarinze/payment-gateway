package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"git.biggorilla.tech/gateway/payment-gateway/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumService interface {
	GenerateNewAddress() *model.Address
	GetBalance(address string) string
}

type ethereumService struct {
	client *ethclient.Client
	ctx    context.Context
}

func NewEthereumService(ctx context.Context) EthereumService {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/5fd8d7c598e4414690cb4f3c49abf585")
	if err != nil {
		log.Fatalf("Oops! There was a problem %x", err)
	} else {
		fmt.Println("Sucess! you are connected to the Ethereum Network")
	}
	return &ethereumService{
		client,
		ctx,
	}
}

func (e *ethereumService) GenerateNewAddress() *model.Address {
	key, _ := crypto.GenerateKey()
	public_key := crypto.PubkeyToAddress(key.PublicKey).Hex()
	private_key := hex.EncodeToString(key.D.Bytes())
	readable_key := &model.Address{
		PublicKey:  public_key,
		PrivateKey: private_key,
	}
	fmt.Println(readable_key)
	return readable_key
}

func (e *ethereumService) GetBalance(address string) string {
	balance, _ := e.client.BalanceAt(e.ctx, common.HexToAddress(address), nil)
	return balance.String()
}
