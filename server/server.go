package server

import (
	"log"
	"net/http"

	h "github.com/akproger/url-screenshot-backend/handlers"    // Алиас h для handlers
	mw "github.com/akproger/url-screenshot-backend/middleware" // Алиас mw для middleware
	"github.com/gorilla/mux"
)

func Start() {
	r := mux.NewRouter()

	// Маршруты для аутентификации
	r.HandleFunc("/api/register", h.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", h.LoginHandler).Methods("POST")

	// CRUD API для URL с защитой миддлвером
	api := r.PathPrefix("/api").Subrouter()
	api.Use(mw.AuthMiddleware)                                                // используем mw.AuthMiddleware
	api.HandleFunc("/urls", h.CreateURLHandler).Methods("POST")               // Создание URL
	api.HandleFunc("/urls/{id:[0-9]+}", h.GetURLHandler).Methods("GET")       // Получение URL по ID
	api.HandleFunc("/urls/{id:[0-9]+}", h.UpdateURLHandler).Methods("PUT")    // Обновление URL
	api.HandleFunc("/urls/{id:[0-9]+}", h.DeleteURLHandler).Methods("DELETE") // Удаление URL

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
