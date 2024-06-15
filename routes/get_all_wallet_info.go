package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/dto/response"
)

func (r *Repository) getAllWalletInfo(ctx *fiber.Ctx) error {
	var wallets []response.Wallets
	if err := r.DB.Select("wallet_id", "balance").Find(&wallets).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Wallets Could Not be Retrieved"})
	}
	walletResponse := response.WalletResponse{Wallets: wallets}
	return ctx.Status(fiber.StatusOK).JSON(walletResponse)
}
