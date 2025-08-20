package main

import (
	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
	"log"
	"os"
)

func main() {
	// Получаем путь к файлу БД из переменной окружения или используем значение по умолчанию
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		// Указываем, что файл БД должен находиться в корне проекта,
		// а не в папке cmd.
		dbFile = "scheduler.db"
	}

	// Инициализируем БД. Эта функция создаст файл и таблицы, если их нет.
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("не удалось инициализировать базу данных: %s", err)
	}
	// Отложенно закрываем соединение с БД при завершении работы программы
	defer db.DB.Close()

	// Запускаем сервер
	server.Run()
}
