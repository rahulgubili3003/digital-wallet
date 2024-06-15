package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/handlers"
	"github.com/rahulgubili3003/digital-wallet/model"
	"io"
	"log"
	"net/http"
	"strings"
)

type ResponseData struct {
	Ok   bool `json:"ok"`
	Data bool `json:"data"`
}

func (r *Repository) CreateWallet(ctx *fiber.Ctx) error {

	fmt.Println("Create Wallet")

	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Missing Authorization Header"})
	}

	splitAuthHeader := strings.Split(authHeader, " ")

	if len(splitAuthHeader) != 2 || splitAuthHeader[0] != "Bearer" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Invalid Authorization Format"})
	}

	token := splitAuthHeader[1]

	post, err := http.Post("http://localhost:3001/api/v1/validate-jwt", "application/json", bytes.NewBuffer([]byte(token)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not validate Token"})
	}

	defer post.Body.Close()

	body, err := io.ReadAll(post.Body)

	if err != nil {
		return err
	}

	var response ResponseData

	err = json.Unmarshal(body, &response)

	if err != nil {
		log.Fatalf("Failed to Unmarshall: %v", err)
	}

	if response.Data == false {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Auth token Invalid"})
	}

	var request handlers.WalletCreateRequest

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Failed to Parse request Body"})
	}
	// Build entity
	entity := model.Wallet{
		UserId:    request.UserId,
		Balance:   0.0,
		IsDeleted: false,
	}

	if err := r.DB.Create(&entity).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not create Wallet"})
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Successfully created"})
}
