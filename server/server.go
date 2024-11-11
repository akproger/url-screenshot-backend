package server

import (
	"log"
	"net/http"

	"github.com/akproger/url-screenshot-backend/handlers"   // импортируем handlers
	"github.com/akproger/url-screenshot-backend/middleware" // импортируем middleware
	"github.com/gorilla/mux"
)

func Start() {
	r := mux.NewRouter()

	// Маршруты аутентификации
	r.HandleFunc("/api/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")

	// CRUD API для URL с защитой миддлвером
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware) // используем middleware.AuthMiddleware для защиты маршрутов
	api.HandleFunc("/urls", handlers.CreateURLHandler).Methods("POST")
	api.HandleFunc("/urls/{id:[0-9]+}", handlers.GetURLHandler).Methods("GET")
	api.HandleFunc("/urls/{id:[0-9]+}", handlers.UpdateURLHandler).Methods("PUT")
	api.HandleFunc("/urls/{id:[0-9]+}", handlers.DeleteURLHandler).Methods("DELETE")

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
