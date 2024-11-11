package handlers

import (
	"encoding/json"
	"net/http"
)

// Пример других функций, которые могут быть в handlers.go, но без дублирования CheckURLHandler
// Например, можно добавить обработчики для других целей
// или оставить его пустым, если все основные обработчики находятся в urls.go

// Пример функции
type HealthCheckResponse struct {
	Status string `json:"status"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{Status: "OK"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
