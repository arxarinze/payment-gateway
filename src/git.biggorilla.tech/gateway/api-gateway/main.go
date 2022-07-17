package main

import (
	// "context"
	// "flag"
	// "log"

	// "git.biggorilla.tech/gateway/payment-gateway/pb"
	"context"
	"flag"
	"log"

	"git.biggorilla.tech/gateway/api-gateway/internal/models"
	"git.biggorilla.tech/gateway/payment-gateway/pb"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/grpc/metadata"
)

func main() {
	app := fiber.New(fiber.Config{
		//Prefork:      true,
		ServerHeader: "BigGorillaApps",
		AppName:      "Api-Gateway",
	})
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Post("/merchant", func(c *fiber.Ctx) error {
		payload := new(models.Request)
		c.BodyParser(&payload)
		// return fiber.NewError(782, "Custom error message")
		addr := flag.String("addr", "localhost:50051", "the address to connect to")
		conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		connection := pb.NewPaymentGatewayServiceClient(conn)
		ctx, cancel := context.WithCancel(context.Background()) //.WithTimeout(context.Background(), time.Minute)
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwidXNlcm5hbWUiOiJ0ZXN0aW5nIiwiaWF0IjoxNTE2MjM5MDIyfQ.-ZWfmCMqmas7sSoU7y8zWwunWUYL7IGShgRw1ykf-84")
		defer cancel()
		r, err := connection.CreateMerchant(ctx, &pb.MerchantRequest{Name: payload.Name, Email: payload.Email})
		if err != nil {
			log.Printf("Greeting: %s", err)
		}
		log.Printf("Greeting: %s", r)
		return c.JSON(r)
		//return c.JSON("Hello, World!")
	})
	app.Listen(":5001")
}
