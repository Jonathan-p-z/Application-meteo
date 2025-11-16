package handlers

import (
	"encoding/json"
	"net/http"

	"weather-app-backend/services"
)

// AlertsHandler gère GET /api/alerts?city=Paris
func AlertsHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		// ville par défaut, tu peux mettre "Paris" ou autre
		city = "Paris"
	}

	alerts, err := services.GetGlobalWeatherAlerts(r.Context(), city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(alerts)
}
