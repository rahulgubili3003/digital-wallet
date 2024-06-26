package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rahulgubili3003/digital-wallet/client"
	"github.com/rahulgubili3003/digital-wallet/model"
	"github.com/rahulgubili3003/digital-wallet/repository"
	"github.com/rahulgubili3003/digital-wallet/routes"
	"log"
	"os"
)

func main() {
	// Initialise the DB Connection
	config := repository.Init()
	db, err := client.ConnectDatabase(config)
	if err != nil {
		log.Fatal("Db Connection Failed")
	}

	if err := model.MigrateWallet(db); err != nil {
		log.Fatal("Could not migrate Wallet DB")
	}

	if err := model.AutoMigrateTransactionsDB(db); err != nil {
		log.Fatal("Could not migrate Transactions DB")
	}

	r := routes.Repository{
		DB: db,
	}
	// Setup web server
	app := fiber.New()
	r.SetupRoutes(app)
	if err := app.Listen(os.Getenv("APP_PORT")); err != nil {
		log.Fatal("Could not start the Server")
	}
}
