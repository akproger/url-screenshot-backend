package server

import (
	"log"
	"net/http"

	h "github.com/akproger/url-screenshot-backend/handlers"    // Алиас h для handlers
	mw "github.com/akproger/url-screenshot-backend/middleware" // Алиас mw для middleware
	"github.com/gorilla/mux"
)

// Middleware для настройки CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Start() {
	r := mux.NewRouter()

	// Добавляем CORS middleware ко всем маршрутам
	r.Use(corsMiddleware)

	// Аутентификация
	r.HandleFunc("/api/register", h.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", h.LoginHandler).Methods("POST")

	// Маршрут для проверки URL
	r.HandleFunc("/api/check-url", h.CheckURLHandler).Methods("POST", "OPTIONS")

	// Маршрут для создания обсуждения
	r.HandleFunc("/api/create-discussion", h.CreateDiscussionHandler).Methods("POST", "OPTIONS")

	// Маршрут для получения обсуждения по gen
	r.HandleFunc("/api/discussions/{gen}", h.GetDiscussionHandler).Methods("GET", "OPTIONS")

	// CRUD API для URL с защитой middleware
	api := r.PathPrefix("/api").Subrouter()
	api.Use(mw.AuthMiddleware)
	api.HandleFunc("/urls", h.CreateDiscussionHandler).Methods("POST") // Здесь указываем корректную функцию
	api.HandleFunc("/urls/{id:[0-9]+}", h.GetURLHandler).Methods("GET")
	api.HandleFunc("/urls/{id:[0-9]+}", h.UpdateURLHandler).Methods("PUT")
	api.HandleFunc("/urls/{id:[0-9]+}", h.DeleteURLHandler).Methods("DELETE")

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
