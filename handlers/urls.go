package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/akproger/url-screenshot-backend/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type URL struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Gen       string    `json:"gen"`
	UserIP    string    `json:"user_ip"`
	CreatedAt time.Time `json:"created_at"`
}

type CheckUrlRequest struct {
	URL string `json:"url"`
}

type CheckUrlResponse struct {
	Exists bool   `json:"exists"`
	Gen    string `json:"gen,omitempty"`
}

type CreateDiscussionRequest struct {
	URL string `json:"url"`
}

type CreateDiscussionResponse struct {
	Gen string `json:"gen"`
}

// CheckURLHandler — проверяет, существует ли URL в базе
func CheckURLHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var gen string
	err := database.DB.QueryRow("SELECT gen FROM urls WHERE url = $1", req.URL).Scan(&gen)
	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode(CheckUrlResponse{Exists: false})
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CheckUrlResponse{Exists: true, Gen: gen})
}

// CreateDiscussionHandler — создает новое обсуждение для URL
func CreateDiscussionHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateDiscussionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	gen := uuid.New().String()[:16] // Генерируем уникальный идентификатор

	userIP := r.RemoteAddr
	createdAt := time.Now()

	// Выполняем вставку в БД
	_, err := database.DB.Exec(
		"INSERT INTO urls (url, gen, user_ip, created_at) VALUES ($1, $2, $3, $4)",
		req.URL, gen, userIP, createdAt,
	)
	if err != nil {
		http.Error(w, "Failed to create discussion", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CreateDiscussionResponse{Gen: gen})
}

// GetURLHandler — получение URL по ID
func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var url URL
	err := database.DB.QueryRow("SELECT id, url, gen, user_ip, created_at FROM urls WHERE id = $1", id).
		Scan(&url.ID, &url.URL, &url.Gen, &url.UserIP, &url.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(url)
}

// UpdateURLHandler — обновление существующего URL
func UpdateURLHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var url URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec("UPDATE urls SET url = $1, gen = $2, user_ip = $3 WHERE id = $4",
		url.URL, url.Gen, url.UserIP, id)
	if err != nil {
		http.Error(w, "Failed to update URL", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("URL updated successfully"))
}

// DeleteURLHandler — удаление URL по ID
func DeleteURLHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := database.DB.Exec("DELETE FROM urls WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete URL", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("URL deleted successfully"))
}
