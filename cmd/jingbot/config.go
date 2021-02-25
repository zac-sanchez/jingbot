package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey string
	BotID  string
	Time   string
	Port   string
}

func getConfig() *Config {

	_ = godotenv.Load(".env")

	if os.Getenv("JINGBOT_TIME") == "" {
		_ = os.Setenv("JINGBOT_TIME", "45 8 * * 1-5")
	}

	return &Config{
		APIKey: os.Getenv("JINGBOT_API_KEY"),
		BotID:  os.Getenv("JINGBOT_USER_ID"),
		Time:   os.Getenv("JINGBOT_TIME"),
		Port: 	os.Getenv("PORT"),
	}
}
