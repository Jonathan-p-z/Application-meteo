package models

// Weather représente une réponse météo riche.
type Weather struct {
	City      string  `json:"city"`
	Country   string  `json:"country,omitempty"`
	Region    string  `json:"region,omitempty"`
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lon,omitempty"`

	// Conditions actuelles
	Temperature      float64 `json:"temperature"`        // temp_c
	FeelsLike        float64 `json:"feels_like"`         // feelslike_c
	Condition        string  `json:"condition"`          // condition.text
	ConditionIconURL string  `json:"condition_icon_url"` // condition.icon (URL)
	Humidity         int     `json:"humidity"`           // humidity (%)
	WindKph          float64 `json:"wind_kph"`           // vent km/h
	WindDegree       int     `json:"wind_degree"`        // angle
	WindDir          string  `json:"wind_dir"`           // N, NE, ...
	PressureMb       float64 `json:"pressure_mb"`        // hPa approx
	VisibilityKm     float64 `json:"visibility_km"`      // vis_km
	UV               float64 `json:"uv"`                 // indice UV
	AirQualityIndex  float64 `json:"air_quality_index"`  // optionnel (si activé dans ton compte)
	Cloud            int     `json:"cloud,omitempty"`    // nébulosité %

	// Prévisions journalières
	ForecastDays []ForecastDay `json:"forecast_days,omitempty"`

	// Prévisions horaires (24h)
	Hourly []ForecastHour `json:"hourly,omitempty"`

	// Informations d'alerte / risques (simplifiées)
	Alerts []WeatherAlert `json:"alerts,omitempty"`
}

// ForecastDay représente la prévision pour un jour.
type ForecastDay struct {
	Date          string  `json:"date"`
	MinTemp       float64 `json:"min_temp"`
	MaxTemp       float64 `json:"max_temp"`
	AvgTemp       float64 `json:"avg_temp"`
	Condition     string  `json:"condition"`
	ConditionIcon string  `json:"condition_icon_url"`

	ChanceOfRain int     `json:"chance_of_rain"` // %
	ChanceOfSnow int     `json:"chance_of_snow"` // %
	RiskThunder  bool    `json:"risk_thunder"`   // basé sur code météo
	WindMaxKph   float64 `json:"wind_max_kph"`
	GustMaxKph   float64 `json:"gust_max_kph"`
	Sunrise      string  `json:"sunrise"`
	Sunset       string  `json:"sunset"`
	MoonPhase    string  `json:"moon_phase"`
}

// ForecastHour représente la prévision pour une heure.
type ForecastHour struct {
	Time         string  `json:"time"` // "2025-01-01 13:00"
	Temp         float64 `json:"temp"`
	Condition    string  `json:"condition"`
	ChanceOfRain int     `json:"chance_of_rain"`
	WindKph      float64 `json:"wind_kph"`
	GustKph      float64 `json:"gust_kph"`
	PressureMb   float64 `json:"pressure_mb"`
	UV           float64 `json:"uv"`
}

// WeatherAlert représente une alerte / risque simplifiée.
type WeatherAlert struct {
	Type     string `json:"type"`     // "orage", "pluie_abondante", "vents_forts", "chaleur"
	Severity string `json:"severity"` // "modéré", "élevé", "extrême"
	Message  string `json:"message"`  // texte explicatif
}
