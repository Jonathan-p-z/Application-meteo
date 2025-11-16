package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"weather-app-backend/config"
	"weather-app-backend/models"
	"weather-app-backend/utils"
)

// WeatherErrorType décrit la catégorie d'erreur métier.
type WeatherErrorType string

const (
	ErrTypeBadRequest WeatherErrorType = "bad_request"
	ErrTypeNotFound   WeatherErrorType = "not_found"
	ErrTypeUpstream   WeatherErrorType = "upstream_error"
	ErrTypeConfig     WeatherErrorType = "config_error"
	ErrTypeDecode     WeatherErrorType = "decode_error"
	ErrTypeUnknown    WeatherErrorType = "unknown_error"
)

// WeatherError est une erreur riche utilisée par le service.
type WeatherError struct {
	Type    WeatherErrorType
	Message string
	Cause   error
}

func (e *WeatherError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *WeatherError) Unwrap() error {
	return e.Cause
}

// helper pour créer un WeatherError
func newWeatherError(t WeatherErrorType, msg string, cause error) *WeatherError {
	return &WeatherError{Type: t, Message: msg, Cause: cause}
}

// GetWeatherForCity récupère conditions + prévisions pour une ville donnée.
func GetWeatherForCity(ctx context.Context, city string) (*models.Weather, error) {
	log.Printf("[weather] incoming request for city=%q\n", city)

	if city == "" {
		return nil, newWeatherError(ErrTypeBadRequest, "paramètre 'city' manquant", nil)
	}

	apiKey, err := config.GetWeatherAPIKey()
	if err != nil {
		werr := newWeatherError(ErrTypeConfig, err.Error(), err)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}

	baseURL := os.Getenv("WEATHER_API_URL")
	if baseURL == "" {
		werr := newWeatherError(ErrTypeConfig, "WEATHER_API_URL is not set", nil)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		werr := newWeatherError(ErrTypeConfig, "WEATHER_API_URL invalide", err)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}

	q := u.Query()
	q.Set("key", apiKey)
	q.Set("q", city)
	q.Set("lang", "fr")
	q.Set("days", "7")
	q.Set("aqi", "yes")
	q.Set("alerts", "no")
	u.RawQuery = q.Encode()

	log.Printf("[weather] calling external API: %s\n", u.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		werr := newWeatherError(ErrTypeUnknown, "impossible de créer la requête HTTP", err)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}

	resp, err := utils.HTTPClient().Do(req)
	if err != nil {
		werr := newWeatherError(ErrTypeUpstream, "échec de l’appel à l’API météo externe", err)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}
	defer resp.Body.Close()

	log.Printf("[weather] external API status=%d\n", resp.StatusCode)

	// Gestion dédiée selon le code HTTP de WeatherAPI
	if resp.StatusCode == http.StatusBadRequest {
		werr := newWeatherError(ErrTypeBadRequest, "la ville demandée est invalide ou mal formée", nil)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}
	if resp.StatusCode == http.StatusNotFound {
		werr := newWeatherError(ErrTypeNotFound, "ville ou ressource météo introuvable", nil)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}
	if resp.StatusCode >= 500 {
		werr := newWeatherError(ErrTypeUpstream, "l’API météo externe rencontre un problème (erreur 5xx)", nil)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}
	if resp.StatusCode != http.StatusOK {
		werr := newWeatherError(ErrTypeUpstream, fmt.Sprintf("réponse inattendue de l’API météo (status %d)", resp.StatusCode), nil)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}

	// Struct pour forecast.json
	var raw struct {
		Location struct {
			Name    string  `json:"name"`
			Region  string  `json:"region"`
			Country string  `json:"country"`
			Lat     float64 `json:"lat"`
			Lon     float64 `json:"lon"`
		} `json:"location"`
		Current struct {
			TempC      float64 `json:"temp_c"`
			FeelsLikeC float64 `json:"feelslike_c"`
			Humidity   int     `json:"humidity"`
			WindKph    float64 `json:"wind_kph"`
			WindDeg    int     `json:"wind_degree"`
			WindDir    string  `json:"wind_dir"`
			PressureMb float64 `json:"pressure_mb"`
			VisKm      float64 `json:"vis_km"`
			UV         float64 `json:"uv"`
			Cloud      int     `json:"cloud"`
			Condition  struct {
				Text string `json:"text"`
				Icon string `json:"icon"`
				Code int    `json:"code"`
			} `json:"condition"`
			AirQuality map[string]float64 `json:"air_quality"` // nécessite aqi=yes et abonnement adéquat
		} `json:"current"`
		Forecast struct {
			Forecastday []struct {
				Date string `json:"date"`
				Day  struct {
					MaxtempC          float64 `json:"maxtemp_c"`
					MintempC          float64 `json:"mintemp_c"`
					AvgtempC          float64 `json:"avgtemp_c"`
					MaxwindKph        float64 `json:"maxwind_kph"`
					TotalprecipMm     float64 `json:"totalprecip_mm"`
					AvgvisKm          float64 `json:"avgvis_km"`
					Avghumidity       float64 `json:"avghumidity"`
					DailyWillItRain   int     `json:"daily_will_it_rain"`
					DailyChanceOfRain int     `json:"daily_chance_of_rain"`
					DailyWillItSnow   int     `json:"daily_will_it_snow"`
					DailyChanceOfSnow int     `json:"daily_chance_of_snow"`
					Condition         struct {
						Text string `json:"text"`
						Icon string `json:"icon"`
						Code int    `json:"code"`
					} `json:"condition"`
					UV float64 `json:"uv"`
				} `json:"day"`
				Hour []struct {
					Time      string  `json:"time"`
					TempC     float64 `json:"temp_c"`
					Condition struct {
						Text string `json:"text"`
						Icon string `json:"icon"`
						Code int    `json:"code"`
					} `json:"condition"`
					ChanceOfRain int     `json:"chance_of_rain"`
					WindKph      float64 `json:"wind_kph"`
					GustKph      float64 `json:"gust_kph"`
					PressureMb   float64 `json:"pressure_mb"`
					UV           float64 `json:"uv"`
				} `json:"hour"`
				Astro struct {
					Sunrise   string `json:"sunrise"`
					Sunset    string `json:"sunset"`
					MoonPhase string `json:"moon_phase"`
				} `json:"astro"`
			} `json:"forecastday"`
		} `json:"forecast"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		werr := newWeatherError(ErrTypeDecode, "impossible de décoder la réponse de l’API météo", err)
		log.Println("[weather] ERROR:", werr)
		return nil, werr
	}

	// Construire la partie "current"
	w := &models.Weather{
		City:        raw.Location.Name,
		Country:     raw.Location.Country,
		Region:      raw.Location.Region,
		Latitude:    raw.Location.Lat,
		Longitude:   raw.Location.Lon,
		Temperature: raw.Current.TempC,
		FeelsLike:   raw.Current.FeelsLikeC,
		Condition:   raw.Current.Condition.Text,
		// l’API renvoie souvent des URLs sans protocole complet, on préfixe en https si besoin
		ConditionIconURL: ensureHTTPSIcon(raw.Current.Condition.Icon),
		Humidity:         raw.Current.Humidity,
		WindKph:          raw.Current.WindKph,
		WindDegree:       raw.Current.WindDeg,
		WindDir:          raw.Current.WindDir,
		PressureMb:       raw.Current.PressureMb,
		VisibilityKm:     raw.Current.VisKm,
		UV:               raw.Current.UV,
		Cloud:            raw.Current.Cloud,
	}

	// Qualité de l’air : AQI global (par ex. "us-epa-index")
	if idx, ok := raw.Current.AirQuality["us-epa-index"]; ok {
		w.AirQualityIndex = idx
	}

	// Prévisions journalières (limitées à 3 jours pour rester lisible)
	for i, d := range raw.Forecast.Forecastday {
		if i >= 3 {
			break
		}
		fd := models.ForecastDay{
			Date:          d.Date,
			MinTemp:       d.Day.MintempC,
			MaxTemp:       d.Day.MaxtempC,
			AvgTemp:       d.Day.AvgtempC,
			Condition:     d.Day.Condition.Text,
			ConditionIcon: ensureHTTPSIcon(d.Day.Condition.Icon),
			ChanceOfRain:  d.Day.DailyChanceOfRain,
			ChanceOfSnow:  d.Day.DailyChanceOfSnow,
			WindMaxKph:    d.Day.MaxwindKph,
			GustMaxKph:    0, // non fourni directement au niveau Day
			Sunrise:       d.Astro.Sunrise,
			Sunset:        d.Astro.Sunset,
			MoonPhase:     d.Astro.MoonPhase,
			RiskThunder:   isThunderRisk(d.Day.Condition.Code),
		}
		w.ForecastDays = append(w.ForecastDays, fd)
	}

	// Prévisions horaires : on prend les 24 prochaines heures si dispo
	if len(raw.Forecast.Forecastday) > 0 {
		for _, h := range raw.Forecast.Forecastday[0].Hour {
			w.Hourly = append(w.Hourly, models.ForecastHour{
				Time:         h.Time,
				Temp:         h.TempC,
				Condition:    h.Condition.Text,
				ChanceOfRain: h.ChanceOfRain,
				WindKph:      h.WindKph,
				GustKph:      h.GustKph,
				PressureMb:   h.PressureMb,
				UV:           h.UV,
			})
		}
	}

	// Dériver quelques alertes/risk simples à partir des valeurs
	w.Alerts = deriveAlerts(w)

	log.Printf("[weather] success city=%q temp=%.1f condition=%q, days=%d, hourly=%d\n",
		w.City, w.Temperature, w.Condition, len(w.ForecastDays), len(w.Hourly))

	return w, nil
}

// ensureHTTPSIcon s’assure que l’URL d’icône est complète.
func ensureHTTPSIcon(icon string) string {
	if icon == "" {
		return ""
	}
	// WeatherAPI renvoie souvent //cdn.weatherapi.com/...
	if icon[:2] == "//" {
		return "https:" + icon
	}
	return icon
}

// isThunderRisk détermine si un code condition correspond à un risque d’orage.
func isThunderRisk(code int) bool {
	// Liste non exhaustive de codes orage (voir docs WeatherAPI)
	switch code {
	case 1087, 1273, 1276, 1279, 1282:
		return true
	default:
		return false
	}
}

// deriveAlerts génère des "alertes" simplifiées à partir des données.
func deriveAlerts(w *models.Weather) []models.WeatherAlert {
	var alerts []models.WeatherAlert

	// Orages violents
	for _, d := range w.ForecastDays {
		if d.RiskThunder {
			alerts = append(alerts, models.WeatherAlert{
				Type:     "orage",
				Severity: "élevé",
				Message:  "Risque d’orage pour la journée " + d.Date,
			})
			break
		}
	}

	// Pluie abondante (si chance de pluie importante)
	for _, d := range w.ForecastDays {
		if d.ChanceOfRain >= 70 {
			alerts = append(alerts, models.WeatherAlert{
				Type:     "pluie_abondante",
				Severity: "modéré",
				Message:  "Probabilité de pluie importante (" + fmt.Sprintf("%d", d.ChanceOfRain) + "%).",
			})
			break
		}
	}

	// Vents forts (rafales approximées via vent max)
	for _, d := range w.ForecastDays {
		if d.WindMaxKph >= 50 {
			alerts = append(alerts, models.WeatherAlert{
				Type:     "vents_forts",
				Severity: "élevé",
				Message:  fmt.Sprintf("Vents forts attendus (jusqu’à %.0f km/h).", d.WindMaxKph),
			})
			break
		}
	}

	// Alerte chaleur (si max > 30°C)
	for _, d := range w.ForecastDays {
		if d.MaxTemp >= 30 {
			alerts = append(alerts, models.WeatherAlert{
				Type:     "chaleur",
				Severity: "élevé",
				Message:  fmt.Sprintf("Épisode de chaleur (max %.1f°C). Pense à bien t’hydrater.", d.MaxTemp),
			})
			break
		}
	}

	return alerts
}
