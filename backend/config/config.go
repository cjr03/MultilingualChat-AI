package config

import (
	"log"
	"os"

	"github.com/subosito/gotenv"
)

type Config struct {
	Port      string
	OpenAIKey string
}

func Load() *Config {
	gotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("OPENAI_API_KEY is not set in environment")
	}

	return &Config{
		Port:      port,
		OpenAIKey: openAIKey,
	}
}
