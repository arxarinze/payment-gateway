package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"git.biggorilla.tech/gateway/payment-gateway/model"
	tether "git.biggorilla.tech/gateway/payment-gateway/tokens"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumService interface {
	GenerateNewAddress() *model.Address
	GetBalance(address string) string
	GetTokenBalance(address string, symbolContractAddress string) string
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

func (e *ethereumService) GetTokenBalance(addr string, symbolContractAddress string) string {
	address := common.HexToAddress(addr)
	contract := common.HexToAddress(symbolContractAddress)
	instance, err1 := tether.NewTether(contract, e.client)
	if err1 != nil {
		log.Fatal(err1)
	}
	bal, _ := instance.BalanceOf(&bind.CallOpts{}, address)
	dec := big.NewFloat(float64(1000000))
	symbol, _ := instance.Symbol(&bind.CallOpts{})
	balance := new(big.Float).Quo(big.NewFloat(0).SetInt(bal), dec)
	return balance.String() + symbol
}
