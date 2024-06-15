package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/client"
	"github.com/rahulgubili3003/digital-wallet/handlers"
	"github.com/rahulgubili3003/digital-wallet/model"
	"github.com/rahulgubili3003/digital-wallet/repository"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

type Repository struct {
	DB *gorm.DB
}

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

type TopUpRequest struct {
	UserId      uint    `json:"user_id"`
	TopUpAmount float64 `json:"top_up_amount"`
}

func (r *Repository) TopUp(ctx *fiber.Ctx) error {
	topUp := TopUpRequest{}
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

	result := r.DB.Model(&wallet).Where("wallet_id =? AND user_id =?", wallet.WalletId, wallet.UserId).Updates(map[string]interface{}{
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

type TransferAmount struct {
	UserId          uint    `json:"user_id"`
	Amount          float64 `json:"amount"`
	RecipientUserId uint    `json:"recipient_user_id"`
}

func (r *Repository) transferAmount(ctx *fiber.Ctx) error {
	var userWallet model.Wallet
	var recipientWallet model.Wallet
	const query = "user_id =?"
	requestBody := &TransferAmount{}

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
	userWalletInfo := r.DB.Where(query, userId).First(&userWallet).Error

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

	recipientWalletInfo := r.DB.Where(query, recipientUserId).First(&recipientWallet).Error

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

	result := r.DB.Model(&userWallet).Where("wallet_id =? AND user_id =?", userWallet.WalletId, userWallet.UserId).Updates(map[string]interface{}{
		"balance":    deductedUserBal,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not start the Money Transfer"})
	}

	existingRecipientBal := recipientWallet.Balance
	updatedRecipientBal := existingRecipientBal + requestBody.Amount

	recipientResult := r.DB.Model(&recipientWallet).Where("wallet_id =? AND user_id =?", recipientWallet.WalletId, recipientWallet.UserId).Updates(map[string]interface{}{
		"balance":    updatedRecipientBal,
		"updated_at": time.Now(),
	})

	if recipientResult.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Could not complete the Money transfer"})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message":      "Wallet Transferred success",
		"sender_id":    userWallet.UserId,
		"recipient_id": recipientWallet.UserId})
}

type Wallets struct {
	WalletId uint    `json:"wallet_id"`
	Balance  float64 `json:"balance"`
}

type WalletResponse struct {
	Wallets []Wallets `json:"wallet_info"`
}

func (r *Repository) getAllWalletInfo(ctx *fiber.Ctx) error {
	var wallets []Wallets
	if err := r.DB.Select("wallet_id", "balance").Find(&wallets).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Wallets Could Not be Retrieved"})
	}
	response := WalletResponse{Wallets: wallets}
	return ctx.Status(fiber.StatusOK).JSON(response)
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

func (r *Repository) setupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/create-wallet", r.CreateWallet)
	api.Post("/top-up", r.TopUp)
	api.Get("/wallet-info/all", r.getAllWalletInfo)
	api.Post("/transfer", r.transferAmount)
}

func main() {
	// Initialise the DB Connection
	config := repository.Init()
	db, err := client.ConnectDatabase(config)
	if err != nil {
		log.Fatal("Db Connection Failed")
	}
	err = model.MigrateWallet(db)
	if err != nil {
		log.Fatal("Could not migrate DB")
	}
	r := Repository{
		DB: db,
	}
	// Setup web server
	app := fiber.New()
	r.setupRoutes(app)
	err = app.Listen(os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatal("Could not start the Server")
	}
}
