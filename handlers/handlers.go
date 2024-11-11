package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthCheckResponse — структура для ответа проверки состояния
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// HealthCheckHandler — обработчик для проверки состояния сервера
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{Status: "OK"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
