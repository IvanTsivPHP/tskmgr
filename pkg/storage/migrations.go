package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

func RunMigrations(direction string) {
	// Загрузка конфигурации из JSON файла
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Формирование строки подключения на основе конфигурации
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%v/%s?sslmode=%s",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName, config.DBSSLMode)

	// Создание пула соединений с базой данных PostgreSQL
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer pool.Close()

	fmt.Println("Успешное подключение к базе данных!")

	// Указание директории с файлами миграций
	migrationsDir := "./migrations"

	// Чтение файлов из директории с миграциями
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatal(err)
	}

	// Фильтрация и сортировка файлов по порядку
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), fmt.Sprintf(".%s.sql", direction)) {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Сортировка файлов по алфавиту для правильного порядка выполнения
	sort.Strings(migrationFiles)

	// Обработка каждого файла миграции
	for _, file := range migrationFiles {
		filePath := filepath.Join(migrationsDir, file)
		fmt.Printf("Выполнение миграции из файла: %s\n", filePath)

		// Чтение содержимого файла
		sqlContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Не удалось прочитать файл %s: %v", filePath, err)
		}

		// Выполнение SQL-запроса
		_, err = pool.Exec(context.Background(), string(sqlContent))
		if err != nil {
			log.Fatalf("Ошибка выполнения SQL из файла %s: %v", filePath, err)
		}

		fmt.Printf("Файл %s успешно выполнен.\n", filePath)
	}
}
