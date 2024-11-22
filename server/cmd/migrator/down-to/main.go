package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sso/internal/config"
	"sso/internal/utils"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// Путь к миграциям
const migrationsPath = "./migrations"

func main() {
	config := config.MustLoad()
	postgresURL := utils.GetPostgresUrl(config)

	var version int64 = -1
	version, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	if version == -1 {
		log.Fatal("Укажите версию")
	}

	fmt.Printf("DOWN TO: %v\n", version)

	// Подключение к базе данных
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Откат миграций
	if err := goose.DownTo(db, migrationsPath, version); err != nil {
		if errors.Is(err, goose.ErrNoMigrations) {
			fmt.Println("Все миграции выполнены")
		}
		log.Fatalf("Ошибка применения миграций: %v", err)
	}
}
