package server

import (
	"go-final-project/pkg/api"
	"log"
	"net/http"
	"os"
)

// Run создает маршрутизатор, инициализирует все маршруты и запускает сервер.
func Run() {
	// Создаем новый мультиплексор (маршрутизатор).
	// Использование собственного мультиплексора вместо глобального http.DefaultServeMux - лучшая практика.
	mux := http.NewServeMux()

	// Регистрируем обработчики API, передавая им наш маршрутизатор.
	api.Init(mux)

	// Настраиваем обработку статических файлов.
	fileServer := http.FileServer(http.Dir("./web"))
	// Регистрируем обработчик для всех путей, которые не были перехвачены API.
	// ServeMux в Go сначала ищет более точные совпадения (например, "/api/tasks"),
	// и только если не находит, использует более общие ("/").
	// Это гарантирует, что запросы к API не будут перехвачены файловым сервером.
	mux.Handle("/", fileServer)

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("Сервер запущен на порту %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("не удалось запустить сервер: %s", err)
	}
}
