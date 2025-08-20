package main

import (
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"log"
	"os"

	"github.com/joho/godotenv" // Импортируем новую библиотеку
)

func main() {
	// Загружаем переменные из .env файла.
	// Это нужно сделать в самом начале, до того как другие части программы их используют.
	if err := godotenv.Load(); err != nil {
		// Если файла .env нет, это не ошибка, просто выводим сообщение.
		// Приложение будет использовать переменные, установленные в системе.
		log.Println("Файл .env не найден, используются системные переменные окружения")
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализируем БД.
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("не удалось инициализировать базу данных: %s", err)
	}

	// Запускаем сервер. Вся логика маршрутизации и запуска инкапсулирована в пакете server.
	server.Run()
}
