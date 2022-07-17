package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	_ "github.com/lib/pq"

	"git.biggorilla.tech/gateway/payment-gateway/database"
	"git.biggorilla.tech/gateway/payment-gateway/helpers"
	"git.biggorilla.tech/gateway/payment-gateway/middleware"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"git.biggorilla.tech/gateway/payment-gateway/repo"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
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
	identity helpers.Identity
}

func (s *server) CreateMerchant(ctx context.Context, in *pb.MerchantRequest) (*pb.GenericResponse, error) {
	auth, _ := grpcMetadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(auth)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbport, user, password, dbname)
	db := database.NewDatabase(ctx).ConnectDatabase(psqlInfo)
	merchantRepo := repo.NewMerchantRepo(ctx, &db)
	data, err1 := merchantRepo.CreateMerchant(ctx, in.GetName(), in.GetEmail(), id)
	out, err2 := json.Marshal(data)
	if err2 != nil || err1 != nil {
		return &pb.GenericResponse{
			Code:    500,
			Message: err2.Error() + err1.Error(),
		}, nil
	}
	defer db.Close()
	return &pb.GenericResponse{
		Code:    200,
		Message: string(out),
	}, nil
}

func (s *server) GenerateLink(ctx context.Context, in *pb.GenerateLinkRequest) (*pb.GenericResponse, error) {
	auth, _ := grpcMetadata.FromIncomingContext(ctx)
	id := s.identity.GetIdentity(auth)
	in.MerchantId = id
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbport, user, password, dbname)
	db := database.NewDatabase(ctx).ConnectDatabase(psqlInfo)
	merchantRepo := repo.NewMerchantRepo(ctx, &db)
	data, err1 := merchantRepo.GenerateLink(ctx, in.GetMerchantId())
	out, err2 := json.Marshal(data)
	if err2 != nil || err1 != nil {
		return &pb.GenericResponse{
			Code:    500,
			Message: err2.Error() + err1.Error(),
		}, nil
	}
	defer db.Close()
	return &pb.GenericResponse{
		Code:    200,
		Message: string(out),
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	jwtManager := helpers.NewJWTManager("testtest123", 15*time.Minute)
	m := middleware.NewMiddleware(context.Background(), *jwtManager).UnaryInterceptor
	s := grpc.NewServer(grpc.UnaryInterceptor(m))
	pb.RegisterPaymentGatewayServiceServer(s, &server{
		identity: helpers.NewIdentity(context.Background()),
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
