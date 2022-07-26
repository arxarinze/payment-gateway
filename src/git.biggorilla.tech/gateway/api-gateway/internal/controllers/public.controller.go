package controllers

import (
	"context"
	"log"

	"git.biggorilla.tech/gateway/api-gateway/internal/models"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PublicController interface {
	GenerateDepositAddress(c *fiber.Ctx) error
	GetMerchantPaymentLink(c *fiber.Ctx) error
}

type publicController struct {
	address *string
}

// GetMerchantPaymentLink implements PublicController
func (p *publicController) GetMerchantPaymentLink(c *fiber.Ctx) error {
	payload := new(models.MerchantPublicRequest)
	payload.PluginID = c.Params("plugin_id")
	conn, err := grpc.Dial(*p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	connection := pb.NewPaymentGatewayServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background()) //.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := connection.GetPublicMerchantInfo(ctx, &pb.MerchantPublicRequest{
		PluginId: payload.PluginID,
	})
	if err != nil {
		log.Printf("Greeting: %s", err)
	}
	log.Printf("Greeting: %s", r)
	return c.JSON(r)
}

// GenerateDepositAddress implements PaymentController
func (p *publicController) GenerateDepositAddress(c *fiber.Ctx) error {
	payload := new(models.DepositAddressRequest)
	c.BodyParser(&payload)
	conn, err := grpc.Dial(*p.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	connection := pb.NewPaymentGatewayServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background()) //.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	r, err := connection.GenerateDepositAddress(ctx, &pb.DepositAddressRequest{
		Cryptosymbol: payload.Cryptosymbol,
		Network:      payload.Network,
		PluginId:     payload.PluginID,
	})
	if err != nil {
		return c.JSON(err.Error())
	}
	log.Printf("Greeting: %s", r)
	return c.JSON(r)
}

func NewPublicController(ctx context.Context, addr *string) PublicController {
	return &publicController{
		address: addr,
	}
}
