package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const defaultPort = "7540"

// Run запускает веб-сервер
func Run() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// Определяем путь к корневой директории проекта
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../..") // Поднимаемся на два уровня вверх из pkg/server

	// Указываем путь к директории web
	webDir := http.Dir(filepath.Join(root, "web"))

	http.Handle("/", http.FileServer(webDir))

	log.Printf("Сервер запущен на порту %s", port)
	log.Printf("Для доступа используйте http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("не удалось запустить сервер: %s\n", err)
	}
}
