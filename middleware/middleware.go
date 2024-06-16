package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Token struct {
	Token string `json:"Authorization"`
}

func Authorization(ctx *fiber.Ctx) (string, error) {

	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return "", ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Missing Authorization Header"})
	}

	splitAuthHeader := strings.Split(authHeader, " ")

	if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "Bearer" {
		return "", ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Invalid Authorization Format"})
	}
	token := splitAuthHeader[1]
	return token, nil
}
