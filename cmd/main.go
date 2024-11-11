package main

import (
	"github.com/akproger/url-screenshot-backend/database"
	"github.com/akproger/url-screenshot-backend/server"
)

func main() {
	// Подключение к базе данных
	database.Connect()

	// Запуск сервера
	server.Start()
}
