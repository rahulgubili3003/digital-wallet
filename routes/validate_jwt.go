package routes

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
)

func (r *Repository) validateJwt(ctx *fiber.Ctx, err error, token string) ([]byte, error, bool) {
	post, err := http.Post("http://localhost:3001/api/v1/validate-jwt", "application/json", bytes.NewBuffer([]byte(token)))
	if err != nil {
		return nil, ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not validate Token"}), true
	}

	defer post.Body.Close()

	body, err := io.ReadAll(post.Body)

	if err != nil {
		return nil, err, true
	}
	return body, nil, false
}
