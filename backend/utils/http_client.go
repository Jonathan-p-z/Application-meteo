package utils

import (
	"net/http"
	"time"
)

var defaultClient = &http.Client{
	Timeout: 5 * time.Second,
}

func HTTPClient() *http.Client {
	return defaultClient
}
