package model

import "gorm.io/gorm"

type Transactions struct {
	gorm.Model
	SenderId          uint    `json:"sender_id"`
	RecipientId       uint    `json:"recipient_id"`
	AmountTransferred float64 `json:"amount_transferred"`
}

func AutoMigrateTransactionsDB(db *gorm.DB) error {
	err := db.AutoMigrate(&Transactions{})
	return err
}
