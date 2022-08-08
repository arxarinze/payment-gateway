package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	tether "git.biggorilla.tech/gateway/payment-gateway/tokens"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
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
				//if common.HexToAddress(vLog.Topics[2].String()).String() == "0x34f53290A60B42DD9D80ECC6b46aB5F4C320144C" {
				fmt.Println("\n", vLog.TxHash.Hex())
				fmt.Println("value", data)
				fmt.Println("from", common.HexToAddress(vLog.Topics[1].String()).String())
				fmt.Println("to", common.HexToAddress(vLog.Topics[2].String()).String())
				//}
			}
		}
	}
}
