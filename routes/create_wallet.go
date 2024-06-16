package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/handlers"
	"github.com/rahulgubili3003/digital-wallet/middleware"
	"github.com/rahulgubili3003/digital-wallet/model"
	"log"
)

func (r *Repository) CreateWallet(ctx *fiber.Ctx) error {

	token, err := middleware.Authorization(ctx)
	if err != nil {
		return err
	}

	body, err, isJwtInvalid := r.validateJwt(ctx, err, token)
	if isJwtInvalid {
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
