package database

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
    "os"
)

var DB *sql.DB

func Connect() {
    connStr := "host=" + os.Getenv("DB_HOST") +
               " user=" + os.Getenv("DB_USER") +
               " password=" + os.Getenv("DB_PASSWORD") +
               " dbname=" + os.Getenv("DB_NAME") +
               " port=" + os.Getenv("DB_PORT") +
               " sslmode=disable"

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Ошибка подключения к БД: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Ошибка проверки соединения с БД: %v", err)
    }

    log.Println("Подключение к БД установлено")
}
