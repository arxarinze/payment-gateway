package main

import (
	"context"
	"flag"
	"git.biggorilla.tech/gateway/api-gateway/internal/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "the address to connect to")
	app := fiber.New(fiber.Config{
		//Prefork:      true,
		ServerHeader: "BigGorillaApps",
		AppName:      "Api-Gateway",
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	authController := controllers.NewAuthController(context.Background())
	merchantController := controllers.NewMerchantController(context.Background(), addr)
	linkController := controllers.NewLinkController(context.Background(), addr)
	publicController := controllers.NewPublicController(context.Background(), addr)
	api := app.Group("/api")
	v1 := api.Group("/v1")
	auth := v1.Group("/auth")

	//authentication
	auth.Post("/login", authController.Login)
	v1.Get("/public/merchant/:plugin_id", publicController.GetMerchantPaymentLink)

	v1.Post("/public/address", publicController.GenerateDepositAddress)

	v1.Get("/link", linkController.GetLink)

	v1.Post("/link", linkController.CreateLink)

	v1.Post("/merchant", merchantController.CreateMerchant)
	app.Listen(":5001")
}
