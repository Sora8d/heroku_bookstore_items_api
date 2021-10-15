package config

import (
	"os"

	"github.com/Sora8d/bookstore_utils-go/logger"
	"github.com/joho/godotenv"
)

type config map[string]string

func init() {
	if err := godotenv.Load("test_envs.env"); err != nil {
		logger.Error("Error loading environment variables", err)
		panic(err)
	}

	Config = config{
		"items_postgres_username": os.Getenv("items_postgres_username"),
		"items_postgres_password": os.Getenv("items_postgres_password"),
		"items_postgres_schema":   os.Getenv("items_postgres_schema"),
		"items_postgres_host":     os.Getenv("items_postgres_host"),
	}
}

var Config config
