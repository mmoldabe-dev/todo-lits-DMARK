package database

import (
	"database/sql"
	"fmt"
	"log"
	"todo-lits-DMARK/app/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Подключение к бд
func Connect() {
	cfg := config.AppConfig
	dsn := fmt.Sprintf(
		"host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	DB = db
	log.Println("Database connected successfully")
}
