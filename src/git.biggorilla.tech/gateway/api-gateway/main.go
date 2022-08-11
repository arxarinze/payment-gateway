package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"git.biggorilla.tech/gateway/api-gateway/internal/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
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
	publicController := controllers.NewPublicController(context.Background(), addr)
	transactionController := controllers.NewTransactionController(context.Background(), addr)
	api := app.Group("/api")
	v1 := api.Group("/v1")
	auth := v1.Group("/auth")
	auth.Post("/login", authController.Login)
	v1.Get("/public/merchant/:plugin_id", publicController.GetMerchantPaymentLink)

	v1.Post("/public/address", publicController.GenerateDepositAddress)
	app.Use(jwtware.New(jwtware.Config{
		Filter: func(ctx *fiber.Ctx) bool {
			headers := ctx.GetReqHeaders()
			split := strings.Split(headers["Authorization"], " ")
			if len(split) <= 1 {
				return false
			}
			token := split[1]
			client := http.Client{}
			req, err := http.NewRequest("GET", "http://ec2-52-72-83-242.compute-1.amazonaws.com/global/verify-sso", nil)
			if err != nil {
				return false
			}

			req.Header = http.Header{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer " + token},
			}

			res, err := client.Do(req)
			if err != nil {
				return false
			}
			var data map[string]interface{}
			json.NewDecoder(res.Body).Decode(&data)
			fmt.Println(data)
			if data["status"] == false {
				return false
			}
			ctx.Locals("user", data)
			ctx.Locals("token", token)
			return true
		},
		SuccessHandler: func(*fiber.Ctx) error {
			fmt.Println("ori")
			return nil
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return fiber.ErrForbidden
		},
		AuthScheme:    "Bearer",
		SigningMethod: "HS256",
		SigningKey:    []byte(""),
	}))
	merchantController := controllers.NewMerchantController(context.Background(), addr)
	linkController := controllers.NewLinkController(context.Background(), addr)

	//authentication

	v1.Get("/link", linkController.GetLink)

	v1.Post("/link", linkController.CreateLink)

	v1.Post("/merchant", merchantController.CreateMerchant)
	v1.Get("/merchant", merchantController.GetMerchants)
	v1.Patch("/merchant", merchantController.UpdateMerchant)
	v1.Get("/transactions", transactionController.GetTransactions)
	app.Listen(":5001")
}
