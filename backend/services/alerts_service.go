package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"weather-app-backend/utils"
)

// ...existing code (si tu as déjà d'autres services)...

type WeatherAlert struct {
	Headline string `json:"headline"`
	Severity string `json:"severity"`
	Areas    string `json:"areas"`
	Event    string `json:"event"`
	Desc     string `json:"desc"`
}

// GetGlobalWeatherAlerts récupère les alertes météo pour une ville (ou zone) donnée.
func GetGlobalWeatherAlerts(ctx context.Context, q string) ([]WeatherAlert, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("weather API key missing (WEATHER_API_KEY)")
	}

	forecastURL := "http://api.weatherapi.com/v1/forecast.json"

	u, err := url.Parse(forecastURL)
	if err != nil {
		return nil, fmt.Errorf("invalid forecast URL: %w", err)
	}

	qp := u.Query()
	qp.Set("key", apiKey)
	qp.Set("q", q)      // ville / coord / code pays, etc.
	qp.Set("days", "1") // un jour suffit pour les alertes
	qp.Set("alerts", "yes")
	u.RawQuery = qp.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating forecast request failed: %w", err)
	}

	resp, err := utils.HTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling forecast API failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("forecast API returned status %d", resp.StatusCode)
	}

	var raw struct {
		Alerts struct {
			Alert []struct {
				Headline string `json:"headline"`
				Severity string `json:"severity"`
				Areas    string `json:"areas"`
				Event    string `json:"event"`
				Desc     string `json:"desc"`
			} `json:"alert"`
		} `json:"alerts"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding forecast API response failed: %w", err)
	}

	alerts := make([]WeatherAlert, 0, len(raw.Alerts.Alert))
	for _, a := range raw.Alerts.Alert {
		alerts = append(alerts, WeatherAlert{
			Headline: a.Headline,
			Severity: a.Severity,
			Areas:    a.Areas,
			Event:    a.Event,
			Desc:     a.Desc,
		})
	}
	return alerts, nil
}
