package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akproger/url-screenshot-backend/database"
	"github.com/gocolly/colly" // для извлечения заголовка из страницы
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type URL struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Gen       string    `json:"gen"`
	UserIP    string    `json:"user_ip"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
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

// ExtractTitle — извлекает заголовок из <h1> или <title>
func ExtractTitle(url string) (string, error) {
	c := colly.NewCollector()
	var title string
	foundH1 := false

	// Проверяем на наличие тега h1
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		if !foundH1 {
			title = e.Text
			foundH1 = true
		}
	})

	// Если h1 не найден, берем из title
	c.OnHTML("title", func(e *colly.HTMLElement) {
		if !foundH1 {
			title = e.Text
		}
	})

	err := c.Visit(url)
	if err != nil {
		return "", fmt.Errorf("не удалось посетить URL: %v", err)
	}

	if title == "" {
		return "", fmt.Errorf("заголовок не найден")
	}
	return title, nil
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
		log.Printf("Ошибка декодирования запроса: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Извлекаем заголовок страницы
	title, err := ExtractTitle(req.URL)
	if err != nil {
		log.Printf("Ошибка извлечения заголовка: %v", err)
		title = "Default Title" // Устанавливаем значение по умолчанию, если заголовок не найден
	}

	gen := uuid.New().String()[:16]
	userIP := r.RemoteAddr
	createdAt := time.Now()

	// Вставляем URL и заголовок в базу данных
	_, err = database.DB.Exec(
		"INSERT INTO urls (url, gen, user_ip, created_at, title) VALUES ($1, $2, $3, $4, $5)",
		req.URL, gen, userIP, createdAt, title,
	)
	if err != nil {
		log.Printf("Ошибка вставки в базу данных: %v", err)
		http.Error(w, "Failed to create discussion", http.StatusInternalServerError)
		return
	}

	log.Println("Обсуждение успешно создано с заголовком:", title)
	json.NewEncoder(w).Encode(CreateDiscussionResponse{Gen: gen})
}

// GetURLHandler — получение URL по ID
func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var url URL
	err := database.DB.QueryRow("SELECT id, url, gen, user_ip, created_at, title FROM urls WHERE id = $1", id).
		Scan(&url.ID, &url.URL, &url.Gen, &url.UserIP, &url.CreatedAt, &url.Title)
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

	_, err := database.DB.Exec("UPDATE urls SET url = $1, gen = $2, user_ip = $3, title = $4 WHERE id = $5",
		url.URL, url.Gen, url.UserIP, url.Title, id)
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

// GetDiscussionHandler — получение обсуждения по генерированному идентификатору
func GetDiscussionHandler(w http.ResponseWriter, r *http.Request) {
	gen := mux.Vars(r)["gen"]

	var url URL
	err := database.DB.QueryRow("SELECT id, url, gen, user_ip, created_at, title FROM urls WHERE gen = $1", gen).
		Scan(&url.ID, &url.URL, &url.Gen, &url.UserIP, &url.CreatedAt, &url.Title)
	if err == sql.ErrNoRows {
		http.Error(w, "Discussion not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(url)
}
