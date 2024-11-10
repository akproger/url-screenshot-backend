package server

import (
	"log"
	"net/http"

	"github.com/akproger/url-screenshot-backend/handlers"
	"github.com/gorilla/mux"
)

func Start() {
	r := mux.NewRouter()

	// Маршруты
	r.HandleFunc("/api/check-url", handlers.CheckURLHandler).Methods("GET")
	// Дополнительные маршруты будут добавлены позже

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
