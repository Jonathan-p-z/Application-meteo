package models

type WeatherResponse struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
