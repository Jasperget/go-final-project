package main

import (
	"go-final-project/pkg/api"
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

	// 1. Инициализируем хранилище (Storage)
	storage, err := db.New(cfg.DBFile)
	if err != nil {
		log.Fatalf("не удалось инициализировать хранилище: %s", err)
	}
	defer storage.Close()

	// 2. Создаем экземпляр API, передавая ему хранилище
	apiHandler := api.New(storage)

	// 3. Запускаем сервер, передавая ему обработчики API
	log.Printf("Сервер запущен на порту %s", cfg.Port)
	if err := server.Run(cfg.Port, apiHandler); err != nil {
		log.Fatalf("не удалось запустить сервер: %s", err)
	}
}
