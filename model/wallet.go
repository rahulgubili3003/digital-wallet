package model

import (
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	WalletId  uint      `gorm:"primary key;autoIncrement" json:"wallet_id"`
	UserId    uint      `gorm:"uniqueIndex:user_id" json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

func MigrateWallet(db *gorm.DB) error {
	err := db.AutoMigrate(&Wallet{})
	return err
}
