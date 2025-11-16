package tests

import (
	"context"
	"testing"

	"weather-app-backend/services"
)

func TestGetWeatherForCity(t *testing.T) {
	ctx := context.Background()
	w, err := services.GetWeatherForCity(ctx, "Paris")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.City != "Paris" {
		t.Fatalf("expected city Paris, got %s", w.City)
	}
}
