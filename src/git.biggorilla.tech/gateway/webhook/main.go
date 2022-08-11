package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	_ "github.com/lib/pq"

	"git.biggorilla.tech/gateway/payment-gateway/database"
	tether "git.biggorilla.tech/gateway/payment-gateway/tokens"
	"git.biggorilla.tech/gateway/webhook/internal/models"
	repo "git.biggorilla.tech/gateway/webhook/internal/repo"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	host     = "ec2-34-229-123-251.compute-1.amazonaws.com"
	dbport   = 5432
	user     = "postgres"
	password = "Moonrider15h3r3t0st0ayFr0st1f0raw1132093@@3340@"
	dbname   = "payment"
)

var (
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbport, user, password, dbname)
)

func main() {
	ctx := context.Background()
	db := database.NewDatabase(ctx).ConnectDatabase(psqlInfo)
	webhookRepo := repo.NewWebhookRepo(ctx, &db)
	defer db.Close()
	client, _ := ethclient.Dial("wss://ropsten.infura.io/ws/v3/5fd8d7c598e4414690cb4f3c49abf585")

	contractAddress := common.HexToAddress("0xB404c51BBC10dcBE948077F18a4B8E553D160084")
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Default().Panicln(err)
	}
	logs := make(chan types.Log)
	a, _ := new(big.Int).SetString(header.Number.String(), 10)
	query := ethereum.FilterQuery{
		FromBlock: a,
		Addresses: []common.Address{
			contractAddress,
		},
	}
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Default().Panicln(err)
	}
	contractAbi, err := abi.JSON(strings.NewReader(string(tether.TetherABI)))
	if err != nil {
		log.Default().Panicln(err)
	}
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	for {
		select {
		case err := <-sub.Err():
			log.Default().Panicln(err)
		case vLog := <-logs:
			if vLog.Topics[0].Hex() == logTransferSigHash.Hex() {
				data, err := contractAbi.Unpack("Transfer", vLog.Data)
				if err != nil {
					log.Default().Panicln(err)
				}
				isAddress, err := webhookRepo.CheckForAddress(ctx, common.HexToAddress(vLog.Topics[2].String()).String())
				if err != nil {
					log.Default().Panicln(err)
				}
				if isAddress {
					value := data[len(data)-1]
					webhookRepo.InsertTransaction(ctx, models.Transaction{
						TxHash: vLog.TxHash.Hex(),
						From:   common.HexToAddress(vLog.Topics[1].String()).Hex(),
						To:     common.HexToAddress(vLog.Topics[2].String()).Hex(),
						Value:  value.(*big.Int).String(),
					})
					fmt.Println("\n", vLog.TxHash.Hex(), " ", vLog.Address.Hex())
					fmt.Println("value", value)
					fmt.Println("from", common.HexToAddress(vLog.Topics[1].String()).Hex())
					fmt.Println("to", common.HexToAddress(vLog.Topics[2].String()).Hex())
				}
			}
		}
	}
}
