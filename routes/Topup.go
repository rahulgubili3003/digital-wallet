package routes

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/constants"
	"github.com/rahulgubili3003/digital-wallet/dto/request"
	"github.com/rahulgubili3003/digital-wallet/model"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (r *Repository) TopUp(ctx *fiber.Ctx) error {
	topUp := request.TopUpRequest{}

	if err := ctx.BodyParser(&topUp); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Failed to Parse the Request"})
	}
	wallet, err := r.findWalletByUserId(topUp.UserId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not find Wallet"})
	}
	existingBal := wallet.Balance
	newBal := existingBal + topUp.TopUpAmount
	result := r.DB.Model(&wallet).Where(constants.WalletAndUserIdQuery, wallet.WalletId, wallet.UserId).Updates(map[string]interface{}{
		"balance":    newBal,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		// Handle error
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update wallet balance",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message":   "Top up successful",
		"wallet_id": wallet.WalletId,
		"balance":   wallet.Balance})
}

func (r *Repository) findWalletByUserId(userId uint) (*model.Wallet, error) {
	var wallet model.Wallet
	result := r.DB.Where("user_id =?", userId).First(&wallet)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// No wallet found for the given user_id
			return nil, fmt.Errorf("no wallet found for user_id %d", userId)
		} else {
			// Some other error occurred
			return nil, result.Error
		}
	}
	return &wallet, nil
}
