package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// DB - глобальная переменная для хранения пула подключений к БД.
var DB *sql.DB

// Init инициализирует соединение с базой данных и при необходимости создает схему.
func Init(dbFile string) error {
	if DB != nil {
		return nil
	}

	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Формируем строку подключения с параметром _busy_timeout.
	// Это заставит соединение ждать до 5 секунд, если база данных заблокирована.
	connStr := fmt.Sprintf("file:%s?_busy_timeout=5000", dbFile)

	// Открываем (или создаем) файл БД
	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(1)

	if install {
		createTableSQL := `CREATE TABLE scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date CHAR(8) NOT NULL DEFAULT "",
            title VARCHAR(255) NOT NULL DEFAULT "",
            comment TEXT NOT NULL DEFAULT "",
            repeat VARCHAR(128) NOT NULL DEFAULT ""
        );`
		if _, err := db.Exec(createTableSQL); err != nil {
			db.Close()
			return err
		}

		createIndexSQL := `CREATE INDEX idx_date ON scheduler (date);`
		if _, err := db.Exec(createIndexSQL); err != nil {
			db.Close()
			return err
		}
	}

	DB = db
	return nil
}

// Close закрывает соединение с базой данных.
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// AddTask добавляет новую задачу в базу данных и возвращает ее ID.
func AddTask(task Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
