package middleware

import (
	"jwt-vi-du-mau/util"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("jwt")

	if _, err := util.ParseJWT(cookie); err != nil {
		ctx.Status(fiber.StatusUnauthorized)
		return ctx.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	return ctx.Next()
}
