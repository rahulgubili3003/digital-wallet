package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/dto/response"
	"github.com/rahulgubili3003/digital-wallet/middleware"
	"log"
)

func (r *Repository) getAllWalletInfo(ctx *fiber.Ctx) error {
	token, err := middleware.Authorization(ctx)
	if err != nil {
		return err
	}
	body, err, isJwtInvalid := r.validateJwt(ctx, err, token)
	if isJwtInvalid {
		return err
	}
	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Fatalf("Failed to Unmarshall: %v", err)
	}
	if responseData.Data == false {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Auth token Invalid"})
	}
	var wallets []response.Wallets
	if err := r.DB.Select("wallet_id", "balance").Find(&wallets).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Wallets Could Not be Retrieved"})
	}
	walletResponse := response.WalletResponse{Wallets: wallets}
	return ctx.Status(fiber.StatusOK).JSON(walletResponse)
}
