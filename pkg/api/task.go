package api

import (
	"database/sql"
	"encoding/json"
	"go-final-project/pkg/db"
	"net/http"
	"time"
)

// getTaskHandler обрабатывает GET-запрос на получение одной задачи.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeError(w, "ошибка получения задачи: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, task, http.StatusOK)
}

// updateTaskHandler обрабатывает PUT-запрос на обновление задачи.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task // Используем db.Task, так как ожидаем поле ID
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация полей (аналогично добавлению задачи)
	if task.Title == "" {
		writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	now := time.Now()
	var taskDate time.Time
	var err error

	if task.Date == "" {
		taskDate = now
	} else {
		// Сначала пытаемся разобрать формат от браузера (DD.MM.YYYY)
		taskDate, err = time.Parse("02.01.2006", task.Date)
		if err != nil {
			// Если не вышло, пытаемся разобрать формат для тестов (YYYYMMDD)
			taskDate, err = time.Parse(DateFormat, task.Date)
			if err != nil {
				writeError(w, "неверный формат даты, ожидается DD.MM.YYYY или YYYYMMDD", http.StatusBadRequest)
				return
			}
		}
	}
	// Преобразуем дату в формат для хранения в БД (YYYYMMDD)
	task.Date = taskDate.Format(DateFormat)

	if task.Repeat != "" {
		nextDateStr, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "некорректное правило повторения: "+err.Error(), http.StatusBadRequest)
			return
		}
		if taskDate.Before(now.Truncate(24 * time.Hour)) {
			task.Date = nextDateStr
		}
	} else {
		if taskDate.Before(now.Truncate(24 * time.Hour)) {
			task.Date = now.Format(DateFormat)
		}
	}

	// Обновление задачи в БД
	err = db.UpdateTask(task)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeError(w, "не удалось обновить задачу: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Отправка успешного ответа
	writeJSON(w, map[string]string{}, http.StatusOK)
}

// deleteTaskHandler обрабатывает DELETE-запрос на удаление задачи.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeError(w, "не удалось удалить задачу: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}
