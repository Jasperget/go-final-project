package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

// DB - глобальная переменная для хранения пула подключений к БД.
var DB *sql.DB

// Init инициализирует соединение с базой данных и при необходимости создает схему.
func Init(dbFile string) error {
	// Проверяем, существует ли файл БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем (или создаем) файл БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Если файл не существовал, создаем в нем таблицы
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

	// Присваиваем соединение глобальной переменной
	DB = db
	return nil
}
