package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const DefaultEnvPath = "./config/.env"

func Load(envPath ...string) error {
	if _, ok := os.LookupEnv("PORT"); ok {
		log.Println("[config] PORT already set in environment, skipping .env load")
		return nil
	}

	path := DefaultEnvPath
	if len(envPath) > 0 && envPath[0] != "" {
		path = envPath[0]
	}

	log.Printf("[config] loading .env from %s\n", path)
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("unable to load env file (%s): %w", path, err)
	}

	log.Println("[config] .env loaded successfully")
	return nil
}

// GetWeatherAPIKey retourne la clé WeatherAPI depuis l'environnement.
func GetWeatherAPIKey() (string, error) {
	key := os.Getenv("WEATHER_API_KEY")
	if key == "" {
		return "", fmt.Errorf("WEATHER_API_KEY is not set")
	}
	return key, nil
}
