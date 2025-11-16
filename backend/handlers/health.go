package handlers

import (
	"encoding/json"
	"net/http"
)

type health struct {
	Status string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(health{Status: "ok"})
}
