package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var template = strings.TrimSpace(`
-- +goose Up
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
`)

const MIGRATIONS_DIR = "./migrations"

func main() {
	var migrationName string
	var fileName string
	var utcDate = getUTCDate()

	flag.StringVar(&migrationName, "migration_name", "", "migration")
	flag.Parse()

	migrationName = "_" + migrationName

	fileName = fmt.Sprintf("%s%s.sql", utcDate, migrationName)

	// Создание файла миграции
	file, err := os.Create(filepath.Join(MIGRATIONS_DIR, fileName))
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}

	// Запись шаблона goose
	_, err = file.WriteString(template)
	if err != nil {
		log.Fatalf("Ошибка записи строки в файл: %v", err)
	}

	file.Close()
	fmt.Printf("File %s successfully created.\n", fileName)
}

func getUTCDate() string {
	currentTime := time.Now().UTC()

	// Форматируем дату и время в формате YYYYMMDDHHMMSS
	return currentTime.Format("20060102150405")
}
