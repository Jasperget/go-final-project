package api

import (
	"database/sql"
	"encoding/json"
	"go-final-project/pkg/dates"
	"go-final-project/pkg/db"
	"net/http"
	"time"
)

// getTaskHandler теперь является методом *API.
func (api *API) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		api.writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := api.storage.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			api.writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			api.writeError(w, "ошибка получения задачи: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	api.writeJSON(w, task, http.StatusOK)
}

// updateTaskHandler теперь является методом *API.
func (api *API) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		api.writeError(w, "ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		api.writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	now := time.Now()
	var taskDate time.Time
	var err error

	if task.Date == "" {
		taskDate = now
	} else {
		taskDate, err = time.Parse(dates.LayoutUser, task.Date)
		if err != nil {
			taskDate, err = time.Parse(dates.LayoutDB, task.Date)
			if err != nil {
				api.writeError(w, "неверный формат даты, ожидается DD.MM.YYYY или YYYYMMDD", http.StatusBadRequest)
				return
			}
		}
	}
	task.Date = taskDate.Format(dates.LayoutDB)

	if task.Repeat != "" {
		nextDateStr, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			api.writeError(w, "некорректное правило повторения: "+err.Error(), http.StatusBadRequest)
			return
		}
		if taskDate.Before(now.Truncate(24 * time.Hour)) {
			task.Date = nextDateStr
		}
	} else {
		if taskDate.Before(now.Truncate(24 * time.Hour)) {
			task.Date = now.Format(dates.LayoutDB)
		}
	}

	err = api.storage.UpdateTask(task)
	if err != nil {
		if err == sql.ErrNoRows {
			api.writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			api.writeError(w, "не удалось обновить задачу: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	api.writeJSON(w, map[string]string{}, http.StatusOK)
}

// deleteTaskHandler теперь является методом *API.
func (api *API) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		api.writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	err := api.storage.DeleteTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			api.writeError(w, "задача не найдена", http.StatusNotFound)
		} else {
			api.writeError(w, "не удалось удалить задачу: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	api.writeJSON(w, map[string]string{}, http.StatusOK)
}
