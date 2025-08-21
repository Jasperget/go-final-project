package api

import (
	"database/sql"
	"go-final-project/pkg/dates" // Импортируем новый пакет
	"go-final-project/pkg/db"
	"net/http"
	"time"
)

// doneTaskHandler обрабатывает POST-запрос на выполнение задачи.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Получаем задачу, чтобы проверить правило повторения
	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeError(w, "ошибка получения задачи: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Если правило повторения не задано, удаляем задачу
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, "не удалось удалить задачу: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		taskDate, err := time.Parse(dates.LayoutDB, task.Date)
		if err != nil {
			writeError(w, "не удалось разобрать текущую дату задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		nextDate, err := NextDate(taskDate, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "не удалось вычислить следующую дату: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Обновляем дату задачи
		err = db.UpdateDate(id, nextDate)
		if err != nil {
			writeError(w, "не удалось обновить дату задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Отправляем успешный ответ
	writeJSON(w, map[string]string{}, http.StatusOK)
}
