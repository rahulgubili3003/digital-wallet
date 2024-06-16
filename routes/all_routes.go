package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type ResponseData struct {
	Ok   bool `json:"ok"`
	Data bool `json:"data"`
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/create-wallet", r.CreateWallet)
	api.Post("/top-up", r.TopUp)
	api.Get("/wallet-info/all", r.getAllWalletInfo)
	api.Post("/transfer", r.transferAmount)
}
