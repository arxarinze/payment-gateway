package controllers

import (
	"context"

	"git.biggorilla.tech/gateway/api-gateway/internal/models"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Login(c *fiber.Ctx) error
}

type authController struct {
}

func NewAuthController(ctx context.Context) AuthController {
	return &authController{}
}

func (a *authController) Login(c *fiber.Ctx) error {
	payload := new(models.LoginRequest)
	c.BodyParser(&payload)

	return c.JSON(payload)
}
