package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Storage инкапсулирует всю логику работы с базой данных.
type Storage struct {
	db *sql.DB
}

// New создает и возвращает новый экземпляр Storage.
// Эта функция заменяет старую Init.
func New(dbFile string) (*Storage, error) {
	// Формируем строку подключения с параметром _busy_timeout, чтобы избежать блокировок.
	connStr := fmt.Sprintf("file:%s?_busy_timeout=5000", dbFile)

	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	schema := `
        CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date CHAR(8) NOT NULL DEFAULT "",
            title VARCHAR(255) NOT NULL DEFAULT "",
            comment TEXT NOT NULL DEFAULT "",
            repeat VARCHAR(128) NOT NULL DEFAULT ""
        );
        CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
    `
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}

	return &Storage{db: db}, nil
}

// Close корректно закрывает соединение с базой данных.
func (s *Storage) Close() error {
	return s.db.Close()
}

// AddTask теперь является методом Storage.
func (s *Storage) AddTask(task Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
