package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rahulgubili3003/digital-wallet/client"
	"github.com/rahulgubili3003/digital-wallet/handlers"
	"github.com/rahulgubili3003/digital-wallet/model"
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

	err := ctx.BodyParser(&request)
	if err != nil {
		err := ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Failed to Parse request Body"})
		if err != nil {
			return err
		}
	}

	entity := model.Wallet{
		UserId:    request.UserId,
		Balance:   0.0,
		IsDeleted: false,
	}

	if err != nil {
		panic(err)
	}

	err = r.DB.Create(&entity).Error

	if err != nil {
		err := ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not create Wallet"})
		if err != nil {
			return err
		}
	}

	err = ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Successfully created"})
	if err != nil {
		return err
	}
	return nil
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
		err := ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not find Wallet"})
		if err != nil {
			return err
		}
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
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to Load Properties")
	}

	config := &client.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

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

	app := fiber.New()
	r.setupRoutes(app)
	err = app.Listen(":3000")
	if err != nil {
		log.Fatal("Could not start the Server")
	}
}
