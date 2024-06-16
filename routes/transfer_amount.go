package routes

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/constants"
	"github.com/rahulgubili3003/digital-wallet/dto/request"
	"github.com/rahulgubili3003/digital-wallet/middleware"
	"github.com/rahulgubili3003/digital-wallet/model"
	"gorm.io/gorm"
	"time"
)

func (r *Repository) transferAmount(ctx *fiber.Ctx) error {
	token, err := middleware.Authorization(ctx)
	if err != nil {
		return err
	}
	_, err, _ = r.validateJwt(ctx, nil, token) // Assuming validation doesn't change based on previous errors
	if err != nil {
		return err
	}

	var requestBody request.TransferAmount
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "Could not parse request Body"})
	}

	userWallet, recipientWallet, err := r.fetchWallets(requestBody.UserId, requestBody.RecipientUserId)
	if err != nil {
		return err
	}

	if err := r.performBalanceUpdate(userWallet, recipientWallet, requestBody.Amount); err != nil {
		return err
	}

	if err := r.createTransaction(userWallet, recipientWallet, requestBody.Amount); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message":      "Wallet Transferred success",
		"sender_id":    userWallet.UserId,
		"recipient_id": recipientWallet.UserId,
	})
}

// Helper functions for fetching wallets, performing balance updates, and creating transactions would go here.
func (r *Repository) fetchWallets(senderID uint, recipientID uint) (*model.Wallet, *model.Wallet, error) {
	var senderWallet, recipientWallet model.Wallet

	// Fetch sender's wallet
	if err := r.DB.Where(constants.UserIdQuery, senderID).First(&senderWallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("wallet not found for sender ID: %s", senderID)
		}
		return nil, nil, err
	}

	// Fetch recipient's wallet
	if err := r.DB.Where(constants.UserIdQuery, recipientID).First(&recipientWallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("wallet not found for recipient ID: %s", recipientID)
		}
		return nil, nil, err
	}

	return &senderWallet, &recipientWallet, nil
}

func (r *Repository) performBalanceUpdate(senderWallet *model.Wallet, recipientWallet *model.Wallet, amount float64) error {
	// Deduct from sender's balance
	newSenderBalance := senderWallet.Balance - amount
	if newSenderBalance < 0 {
		return fmt.Errorf("insufficient funds")
	}

	// Add to recipient's balance
	newRecipientBalance := recipientWallet.Balance + amount

	// Update balances
	if err := r.DB.Model(&senderWallet).Where(constants.WalletAndUserIdQuery, senderWallet.WalletId, senderWallet.UserId).Updates(map[string]interface{}{
		"balance":    newSenderBalance,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	if err := r.DB.Model(&recipientWallet).Where(constants.WalletAndUserIdQuery, recipientWallet.WalletId, recipientWallet.UserId).Updates(map[string]interface{}{
		"balance":    newRecipientBalance,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) createTransaction(senderWallet *model.Wallet, recipientWallet *model.Wallet, amount float64) error {
	transaction := model.Transactions{
		SenderId:          senderWallet.UserId,
		RecipientId:       recipientWallet.UserId,
		AmountTransferred: amount}
	if err := r.DB.Create(&transaction).Error; err != nil {
		return err
	}
	return nil
}
