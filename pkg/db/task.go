package db

import (
	"database/sql"
	"go-final-project/pkg/dates"
	"time"
)

// Task представляет собой задачу в планировщике.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Tasks теперь является методом Storage.
func (s *Storage) Tasks(limit int, search string) ([]Task, error) {
	var rows *sql.Rows
	var err error

	// Проверяем, является ли строка поиска датой в формате "DD.MM.YYYY"
	t, err := time.Parse(dates.LayoutUser, search)
	if err == nil {
		// Поиск по дате
		dateStr := t.Format(dates.LayoutDB) // Используем константу
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
		rows, err = s.db.Query(query, dateStr, limit)
	} else if search != "" {
		// Поиск по тексту
		searchPattern := "%" + search + "%"
		query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
		rows, err = s.db.Query(query, searchPattern, searchPattern, limit)
	} else {
		// Нет поиска, получаем все задачи
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = s.db.Query(query, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		// Сканируем данные в поля структуры Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Если задач нет, возвращаем пустой срез, а не nil, для корректного JSON-ответа {"tasks":[]}.
	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

// GetTask теперь является методом Storage.
func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	// QueryRow ожидает ровно одну строку. Если задача не найдена, вернется ошибка sql.ErrNoRows.
	err := s.db.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// UpdateTask теперь является методом Storage.
func (s *Storage) UpdateTask(task Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	// Проверяем, была ли обновлена хотя бы одна строка.
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		// Если ни одна строка не была затронута, значит задачи с таким ID не существует.
		return sql.ErrNoRows
	}
	return nil
}

// UpdateDate теперь является методом Storage.
func (s *Storage) UpdateDate(id, newDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := s.db.Exec(query, newDate, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// DeleteTask теперь является методом Storage.
func (s *Storage) DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows // Задача с таким ID не найдена
	}
	return nil
}
