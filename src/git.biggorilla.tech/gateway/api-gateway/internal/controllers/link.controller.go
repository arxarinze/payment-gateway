package controllers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"git.biggorilla.tech/gateway/api-gateway/internal/models"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type LinkController interface {
	GetLink(c *fiber.Ctx) error
	CreateLink(c *fiber.Ctx) error
}

type linkController struct {
	address *string
}

// CreateLink implements LinkController
func (l *linkController) CreateLink(c *fiber.Ctx) error {
	payload := new(models.GenerateLinkRequest)
	c.BodyParser(&payload)
	fmt.Println(payload)
	// return fiber.NewError(782, "Custom error message")
	conn, err := grpc.Dial(*l.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	connection := pb.NewPaymentGatewayServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background()) //.WithTimeout(context.Background(), time.Minute)
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.Locals("token").(string))
	defer cancel()
	r, err := connection.GenerateLink(ctx, &pb.GenerateLinkRequest{MerchantId: payload.MerchantID})
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "duplicate key") {
			return c.JSON(map[string]string{
				"Code":    "409",
				"Message": "Link Already Generated",
			})
		}
		return c.JSON(map[string]string{
			"Code":    "500",
			"Message": "An Internal Error Occured",
		})
	}
	return c.JSON(r)
}

// GetLink implements LinkController
func (l *linkController) GetLink(c *fiber.Ctx) error {
	payload := new(models.GetLinkRequest)
	c.BodyParser(&payload)
	conn, err := grpc.Dial(*l.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	connection := pb.NewPaymentGatewayServiceClient(conn)
	ctx, cancel := context.WithCancel(context.Background()) //.WithTimeout(context.Background(), time.Minute)
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.Locals("token").(string))
	defer cancel()
	r, err := connection.GetPluginLink(ctx, &pb.PluginLinkRequest{MerchantId: payload.MerchantID, Type: payload.Type})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.JSON(map[string]string{
				"Code":    "409",
				"Message": "Link Already Generated",
			})
		}
		return c.JSON(map[string]string{
			"Code":    "500",
			"Message": "An Internal Error Occured",
		})
	}
	return c.JSON(r)
}

func NewLinkController(ctx context.Context, addr *string) LinkController {
	return &linkController{
		address: addr,
	}
}
