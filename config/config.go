package config

import (
	"os"
)

type config map[string]string

func init() {
	Config = config{
		"database": os.Getenv("DATABASE_URL"),
		"address":  os.Getenv("adress"),
		"oauth":    os.Getenv("oauth"),
	}
}

var Config config
