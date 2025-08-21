package main

import (
	"go-final-project/pkg/config" // Импортируем новый пакет
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные из .env файла.
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются системные переменные окружения")
	}

	// Считываем конфигурацию ОДИН РАЗ при старте приложения.
	if err := config.Load(); err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %s", err)
	}
	cfg := config.Get()

	// Инициализируем БД, используя значение из конфигурации.
	if err := db.Init(cfg.DBFile); err != nil {
		log.Fatalf("не удалось инициализировать базу данных: %s", err)
	}
	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Printf("Ошибка при закрытии соединения с БД: %v", err)
		}
	}()

	// Запускаем сервер.
	server.Run()
}
