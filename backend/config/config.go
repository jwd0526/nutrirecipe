package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	USDAAPIKey  string
	Port        string
	OllamaURL   string
	OllamaModel string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	return &Config{
		DatabaseURL: require("DATABASE_URL"),
		USDAAPIKey:  require("USDA_API_KEY"),
		Port:        getOrDefault("PORT", "8080"),
		OllamaURL:   getOrDefault("OLLAMA_URL", "http://127.0.0.1:11434"),
		OllamaModel: getOrDefault("OLLAMA_MODEL", "qwen2.5-coder:7b"),
	}
}

func require(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

func getOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
