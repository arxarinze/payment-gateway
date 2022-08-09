package controllers

import (
	"context"
	"log"

	"git.biggorilla.tech/gateway/api-gateway/internal/models"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type TransactionController interface {
	GetTransactions(c *fiber.Ctx) error
}

type transactionController struct {
	address *string
}

// GetTransactions implements TransactionController
func (m *transactionController) GetTransactions(c *fiber.Ctx) error {
	payload := new(models.TransactionRequest)
	c.BodyParser(&payload)
	conn, err := grpc.Dial(*m.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	connection := pb.NewPaymentGatewayServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.Locals("token").(string))
	defer cancel()
	r, err := connection.GetTransactions(ctx, &pb.TransactionRequest{MerchantId: payload.MerchantID})
	if err != nil {
		log.Printf("Greeting: %s", err)
	}

	return c.JSON(r)
}

func NewTransactionController(ctx context.Context, addr *string) TransactionController {
	return &transactionController{
		address: addr,
	}
}
