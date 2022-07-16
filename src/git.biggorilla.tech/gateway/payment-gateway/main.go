package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnsafePaymentGatewayServiceServer
}

func (s *server) CreateMerchant(ctx context.Context, in *pb.MerchantRequest) (*pb.GenericReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.GenericReply{Message: "Created Merchant with name " + in.GetName() + " successfully", Code: 200}, nil
}

func (s *server) GenerateLink(ctx context.Context, in *pb.GenerateLinkRequest) (*pb.GenericReply, error) {
	log.Printf("Received: %v", in.GetMerchantId())
	return &pb.GenericReply{Message: "Created Merchant with name " + in.GetMerchantId() + " successfully", Code: 200}, nil
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
