package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/akproger/url-screenshot-backend/database"
)

type URLCheckResponse struct {
	Exists bool   `json:"exists"`
	URL    string `json:"url,omitempty"`
}

func CheckURLHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM urls WHERE url=$1)", url).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := URLCheckResponse{Exists: exists}
	if exists {
		var existingURL string
		database.DB.QueryRow("SELECT unique_page_url FROM urls WHERE url=$1", url).Scan(&existingURL)
		response.URL = existingURL
	}

	w.Header().Set("Content-Type", application/json)
	json.NewEncoder(w).Encode(response)
}
