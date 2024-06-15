package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/handlers"
	"github.com/rahulgubili3003/digital-wallet/model"
)

func (r *Repository) CreateWallet(ctx *fiber.Ctx) error {
	fmt.Println("Create Wallet")
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
