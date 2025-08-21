package api

import (
	"database/sql"
	"go-final-project/pkg/dates"
	"net/http"
	"time"
)

// doneTaskHandler теперь является методом *API.
func (api *API) doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		api.writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Получаем задачу, чтобы проверить правило повторения
	task, err := api.storage.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			api.writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			api.writeError(w, "ошибка получения задачи: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Если правило повторения не задано, удаляем задачу
	if task.Repeat == "" {
		err = api.storage.DeleteTask(id)
		if err != nil {
			api.writeError(w, "не удалось удалить задачу: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		taskDate, err := time.Parse(dates.LayoutDB, task.Date)
		if err != nil {
			api.writeError(w, "не удалось разобрать текущую дату задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}

		nextDate, err := NextDate(taskDate, task.Date, task.Repeat)
		if err != nil {
			api.writeError(w, "не удалось вычислить следующую дату: "+err.Error(), http.StatusInternalServerError)
			return
		}
		// Обновляем дату задачи
		err = api.storage.UpdateDate(id, nextDate)
		if err != nil {
			api.writeError(w, "не удалось обновить дату задачи: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Отправляем успешный ответ
	api.writeJSON(w, map[string]string{}, http.StatusOK)
}
