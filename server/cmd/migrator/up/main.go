package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sso/internal/config"
	"sso/internal/utils"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// Путь к миграциям
const migrationsPath = "./migrations"

func main() {
	config := config.MustLoad()
	postgresURL := utils.GetPostgresUrl(config)
	// Подключение к базе данных
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Применение миграций
	if err := goose.Up(db, migrationsPath); err != nil {
		if errors.Is(err, goose.ErrNoMigrations) {
			fmt.Println("Все миграции выполнены")
		}
		log.Fatalf("Ошибка применения миграций: %v", err)
	}
}
