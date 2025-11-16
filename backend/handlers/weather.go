package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"weather-app-backend/models"
	"weather-app-backend/services"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "Le paramètre 'city' est obligatoire.",
		})
		return
	}

	data, err := services.GetWeatherForCity(r.Context(), city)
	if err != nil {
		var werr *services.WeatherError
		if errors.As(err, &werr) {
			switch werr.Type {
			case services.ErrTypeBadRequest:
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(models.ErrorResponse{
					Error: "La ville saisie est invalide ou non supportée par l’API météo.",
				})
				return
			case services.ErrTypeNotFound:
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(models.ErrorResponse{
					Error: "Aucune donnée météo trouvée pour cette ville.",
				})
				return
			case services.ErrTypeConfig:
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.ErrorResponse{
					Error: "Erreur de configuration côté serveur (clé API ou URL manquante).",
				})
				return
			case services.ErrTypeUpstream:
				w.WriteHeader(http.StatusBadGateway)
				_ = json.NewEncoder(w).Encode(models.ErrorResponse{
					Error: "L’API météo externe ne répond pas correctement. Réessaie plus tard.",
				})
				return
			case services.ErrTypeDecode:
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.ErrorResponse{
					Error: "Le serveur n’a pas réussi à comprendre la réponse de l’API météo.",
				})
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(models.ErrorResponse{
					Error: "Une erreur interne est survenue lors de la récupération de la météo.",
				})
				return
			}
		}

		// Si ce n’est pas un WeatherError (cas improbable), fallback générique
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "Erreur inattendue côté serveur.",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data) // renvoie *models.Weather complet
}
