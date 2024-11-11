package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/akproger/url-screenshot-backend/database"
	"github.com/gorilla/mux"
)

type URL struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Gen       string `json:"gen"`
	UserIP    string `json:"user_ip"`
	CreatedAt string `json:"created_at"`
}

// CreateURLHandler — добавление нового URL
// CreateURLHandler — добавление нового URL
func CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	var url URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Выполняем вставку и возвращаем ID новой записи
	var lastInsertID int
	err = database.DB.QueryRow("INSERT INTO urls (url, gen, user_ip) VALUES ($1, $2, $3) RETURNING id",
		url.URL, url.Gen, url.UserIP).Scan(&lastInsertID)
	if err != nil {
		http.Error(w, "Failed to create URL", http.StatusInternalServerError)
		return
	}

	// Формируем JSON-ответ с ID новой записи
	response := map[string]interface{}{
		"message": "URL created successfully",
		"id":      lastInsertID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetURLHandler — получение URL по ID
func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	var url URL
	err = database.DB.QueryRow("SELECT id, url, gen, user_ip, created_at FROM urls WHERE id = $1", id).
		Scan(&url.ID, &url.URL, &url.Gen, &url.UserIP, &url.CreatedAt)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(url)
}

// UpdateURLHandler — обновление существующего URL
func UpdateURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	var url URL
	err = json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("UPDATE urls SET url = $1, gen = $2, user_ip = $3 WHERE id = $4",
		url.URL, url.Gen, url.UserIP, id)
	if err != nil {
		http.Error(w, "Failed to update URL", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("URL updated successfully"))
}

// DeleteURLHandler — удаление URL по ID
func DeleteURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM urls WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("URL deleted successfully"))
}
