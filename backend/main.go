package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"weather-app-backend/config"
	"weather-app-backend/handlers"
)

func main() {
	// Charger la config (.env)
	if err := config.Load(); err != nil {
		log.Println("warning: could not load config:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", handlers.HealthHandler)
	mux.HandleFunc("/api/weather", handlers.WeatherHandler)
	mux.HandleFunc("/api/alerts", handlers.AlertsHandler)

	// Page d'accueil + assets front
	// On part du dossier backend et on remonte vers ../frontend
	frontendDir := filepath.Join("..", "frontend")
	fileServer := http.FileServer(http.Dir(frontendDir))

	// Quand on va sur "/", on sert index.html du frontend
	mux.Handle("/", fileServer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server listening on http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
