package main

import (
	"context"
	sql "database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net"

	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"git.biggorilla.tech/gateway/payment-gateway/repo"
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
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnsafePaymentGatewayServiceServer
	merchantRepo repo.MerchantRepo
}

func (s *server) CreateMerchant(ctx context.Context, in *pb.MerchantRequest) (*pb.MerchantRepsonse, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbport, user, password, dbname)
	db := _connectDatbase(psqlInfo)
	s.merchantRepo = repo.NewMerchantRepo(ctx, &db)
	data, err := s.merchantRepo.CreateMerchant(ctx, in.GetName(), in.GetEmail())
	if err != nil {
		panic(err)
	}
	fmt.Println(in, data.Email)
	defer db.Close()
	return &pb.MerchantRepsonse{
		Id:    data.ID,
		Name:  data.Name,
		Email: data.Email,
	}, nil
}

func (s *server) GenerateLink(ctx context.Context, in *pb.GenerateLinkRequest) (*pb.GenericReply, error) {
	log.Printf("Received: %v", in.GetMerchantId())
	return &pb.GenericReply{Message: "Created Merchant with name " + in.GetMerchantId() + " successfully", Code: 200}, nil
}

func _connectDatbase(psqlInfo string) sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return *db
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPaymentGatewayServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
