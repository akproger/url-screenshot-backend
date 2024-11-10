package main

import (
	"url-screenshot-backend/database"
	"url-screenshot-backend/server"
)

func main() {
	// Подключение к базе данных
	database.Connect()

	// Запуск сервера
	server.Start()
}
