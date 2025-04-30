package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"talk/internal/config"

	_ "database/sql"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func MustConnectPostgres(config config.Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", GetPostgresConnectionString(config))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func Migrate(db *sqlx.DB) {
	fullMigrationsPath := path.Join(GetWdPath(), "./migrations")

	// Применение миграций
	if err := goose.Up(db.DB, fullMigrationsPath); err != nil {
		if errors.Is(err, goose.ErrNoMigrations) {
			fmt.Println("Все миграции выполнены")
		}
		log.Fatalf("Ошибка применения миграций: %v", err)
	}
}

func GetWdPath() string {
	wdFromEnv := os.Getenv("WORKDIR_PATH")
	wdFromOs, _ := os.Getwd()

	if wdFromEnv != "" {
		return wdFromEnv
	}

	return wdFromOs
}

func GetPostgresUrl(config config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=disable",
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Database,
	)
}

func GetPostgresConnectionString(config config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Postgres.Host,
		strconv.Itoa(config.Postgres.Port),
		config.Postgres.User,
		config.Postgres.Database,
		config.Postgres.Password,
		"disable",
	)
}

func MergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	// Создаем новую карту, в которую будем сливать
	merged := make(map[string]interface{})

	// Добавляем все элементы из первой карты
	for k, v := range map1 {
		merged[k] = v
	}

	// Добавляем все элементы из второй карты
	for k, v := range map2 {
		merged[k] = v
	}

	return merged
}
