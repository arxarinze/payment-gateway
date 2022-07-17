package main

import (
	"context"
	"flag"
	"log"
	_ "time"

	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcMetadata "google.golang.org/grpc/metadata"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPaymentGatewayServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithCancel(context.Background()) //.WithTimeout(context.Background(), time.Minute)
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwidXNlcm5hbWUiOiJ0ZXN0aW5nIiwiaWF0IjoxNTE2MjM5MDIyfQ.-ZWfmCMqmas7sSoU7y8zWwunWUYL7IGShgRw1ykf-84")
	defer cancel()
	r, err := c.CreateMerchant(ctx, &pb.MerchantRequest{Name: "tnk", Email: "rxarinze@live.com"})
	if err != nil {
		log.Printf("Greeting: %s", err)
	}
	log.Printf("Greeting: %s", r)
}
