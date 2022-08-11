package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	_ "github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"

	"git.biggorilla.tech/gateway/payment-gateway/database"
	"git.biggorilla.tech/gateway/payment-gateway/helpers"
	"git.biggorilla.tech/gateway/payment-gateway/middleware"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"git.biggorilla.tech/gateway/payment-gateway/repo"
	pRPC "git.biggorilla.tech/gateway/payment-gateway/rpc"
	services "git.biggorilla.tech/gateway/payment-gateway/services/web3"
	"google.golang.org/grpc"
)

const (
	host     = "ec2-34-229-123-251.compute-1.amazonaws.com"
	dbport   = 5432
	user     = "postgres"
	password = "Moonrider15h3r3t0st0ayFr0st1f0raw1132093@@3340@"
	dbname   = "payment"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbport, user, password, dbname)
)

func main() {
	flag.Parse()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ctx := context.Background()
	jwtManager := helpers.NewJWTManager("thegfat13224234jahdkskAADJAKDJKAkjskdajdkasj@@@11111jdsdajksk!!!!!!!$$$%%#@", 15*time.Minute)
	middleware := middleware.NewMiddleware(ctx, *jwtManager).UnaryInterceptor
	server := grpc.NewServer(grpc.UnaryInterceptor(middleware))
	db := database.NewDatabase(ctx).ConnectDatabase(psqlInfo)
	identity := helpers.NewIdentity(ctx)
	merchantRepo := repo.NewMerchantRepo(ctx, &db)
	ethereumClient := services.NewEthereumService(ctx)
	//client, _ := ethclient.Dial("wss://mainnet.infura.io/ws/v3/5fd8d7c598e4414690cb4f3c49abf585")
	bala := ethereumClient.GetTokenBalance("0x55FE002aefF02F77364de339a1292923A15844B8", "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	fmt.Println(bala)
	rpc := pRPC.NewRPCInterface(ctx, identity, merchantRepo, ethereumClient)
	pb.RegisterPaymentGatewayServiceServer(server, rpc)
	log.Printf("server listening at %v", listen.Addr())
	defer db.Close()
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
