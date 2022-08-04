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

type MerchantController interface {
	CreateMerchant(c *fiber.Ctx) error
}

type merchantController struct {
	address *string
}

// CreateMerchant implements MerchantController
func (m *merchantController) CreateMerchant(c *fiber.Ctx) error {
	payload := new(models.MerchantRequest)
	c.BodyParser(&payload)
	conn, err := grpc.Dial(*m.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	connection := pb.NewPaymentGatewayServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background()) //.WithTimeout(context.Background(), time.Minute)
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.Locals("token").(string))
	defer cancel()
	r, err := connection.CreateMerchant(ctx, &pb.MerchantRequest{Name: payload.Name, Email: payload.Email})
	if err != nil {
		log.Printf("Greeting: %s", err)
	}
	log.Printf("Greeting: %s", r)
	return c.JSON(r)
}

func NewMerchantController(ctx context.Context, addr *string) MerchantController {
	return &merchantController{
		address: addr,
	}
}
