package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/constants"
	"github.com/rahulgubili3003/digital-wallet/dto/request"
	"github.com/rahulgubili3003/digital-wallet/model"
	"gorm.io/gorm"
	"log"
	"time"
)

func (r *Repository) transferAmount(ctx *fiber.Ctx) error {
	var userWallet model.Wallet
	var recipientWallet model.Wallet

	requestBody := &request.TransferAmount{}

	if err := ctx.BodyParser(requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not parse request Body"})
	}

	userId := requestBody.UserId
	amount := requestBody.Amount
	recipientUserId := requestBody.RecipientUserId

	if amount <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Transfer Amount is invalid. Only positive amount is valid."})
	}
	userWalletInfo := r.DB.Where(constants.UserIdQuery, userId).First(&userWallet).Error

	if userWalletInfo != nil {
		if errors.Is(userWalletInfo, gorm.ErrRecordNotFound) {
			// No wallet found for the given user_id
			return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"message":   "Sender user Id is invalid",
				"sender_id": userId,
			})
		} else {
			// Some other error occurred
			log.Printf("Error Occurred while retrieving wallet Info :%s", userWalletInfo)
			return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"message": "Internal Error",
				"reason":  "Failed to fetch Wallet Details"})
		}
	}

	recipientWalletInfo := r.DB.Where(constants.UserIdQuery, recipientUserId).First(&recipientWallet).Error

	if recipientWalletInfo != nil {
		if errors.Is(recipientWalletInfo, gorm.ErrRecordNotFound) {
			// No wallet found for the given user_id
			return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
				"message":   "Recipient user Id is invalid",
				"sender_id": recipientUserId,
			})
		} else {
			// Some other error occurred
			log.Printf("Error Occurred while retrieving wallet Info :%s", recipientWalletInfo)
			return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"message": "Internal Error",
				"reason":  "Failed to fetch Wallet Details"})
		}
	}

	existingUserBal := userWallet.Balance
	deductedUserBal := existingUserBal - requestBody.Amount

	if deductedUserBal < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Sender Balance insufficient"})
	}

	result := r.DB.Model(&userWallet).Where(constants.WalletAndUserIdQuery, userWallet.WalletId, userWallet.UserId).Updates(map[string]interface{}{
		"balance":    deductedUserBal,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not start the Money Transfer"})
	}

	existingRecipientBal := recipientWallet.Balance
	updatedRecipientBal := existingRecipientBal + requestBody.Amount

	recipientResult := r.DB.Model(&recipientWallet).Where(constants.WalletAndUserIdQuery, recipientWallet.WalletId, recipientWallet.UserId).Updates(map[string]interface{}{
		"balance":    updatedRecipientBal,
		"updated_at": time.Now(),
	})

	if recipientResult.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not complete the Money transfer"})
	}

	transactionEntity := model.Transactions{
		SenderId:          userWallet.UserId,
		RecipientId:       recipientWallet.UserId,
		AmountTransferred: amount,
	}

	if err := r.DB.Create(&transactionEntity).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not register the Transaction Record in the DB"})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message":      "Wallet Transferred success",
		"sender_id":    userWallet.UserId,
		"recipient_id": recipientWallet.UserId})
}
