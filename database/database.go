package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	// Считываем переменные окружения
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Проверяем наличие каждой переменной, чтобы избежать ошибок
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		log.Fatalf("Одно или несколько переменных окружения для подключения к БД не заданы")
	}

	// Формируем строку подключения
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Подключаемся к базе данных
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// Проверяем соединение
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Ошибка проверки соединения с БД: %v", err)
	}

	log.Println("Подключение к БД установлено")
}
