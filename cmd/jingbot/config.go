package main

import (
	"os"
	"strconv"
)

type Config struct {
	APIKey string
	Environment string
	BotID  string
	Minutes int
	Time   string
	Port   string
}

func getConfig() *Config {

	if os.Getenv("JINGBOT_TIME") == "" {
		_ = os.Setenv("JINGBOT_TIME", "45 21 * * 1-5")
	}

	minutes, err := strconv.Atoi(os.Getenv("N_RANDOM_MINUTES"))
	if err != nil {
		panic(err)
	}

	cfg := &Config{
		Time:   os.Getenv("JINGBOT_TIME"),
		Port: 	os.Getenv("PORT"),
		Environment: os.Getenv("JINGBOT_ENVIRONMENT"),
		Minutes: minutes,
	}

	if cfg.Environment == "prod" {
		cfg.APIKey = os.Getenv("JINGBOT_PROD_API_KEY")
		cfg.BotID =  os.Getenv("JINGBOT_PROD_USER_ID")
	} else {
		cfg.APIKey = os.Getenv("JINGBOT_DEV_API_KEY")
		cfg.BotID =  os.Getenv("JINGBOT_DEV_USER_ID")
	}

	return cfg
}
