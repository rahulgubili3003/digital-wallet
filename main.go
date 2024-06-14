package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rahulgubili3003/digital-wallet/client"
	"github.com/rahulgubili3003/digital-wallet/handlers"
	"github.com/rahulgubili3003/digital-wallet/model"
	"gorm.io/gorm"
	"log"
	"os"
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
			"message": "Could not create Book"})
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
	api := app.Group("/api/v1")
	api.Post("/create-wallet", r.CreateWallet)

	err = app.Listen(":3000")
	if err != nil {
		log.Fatal("Could not start the Server")
	}
}
