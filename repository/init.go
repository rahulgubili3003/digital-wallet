package repository

import (
	"github.com/joho/godotenv"
	"github.com/rahulgubili3003/digital-wallet/client"
	"log"
	"os"
)

func Init() *client.Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to Load Properties")
	}

	return &client.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}
