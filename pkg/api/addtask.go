package api

import (
	"encoding/json"
	"go-final-project/pkg/dates"
	"go-final-project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

// TaskPayload — это структура для десериализации JSON при добавлении новой задачи.
type TaskPayload struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// addTaskHandler теперь является методом *API.
func (api *API) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var payload TaskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.writeError(w, "ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if payload.Title == "" {
		api.writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	now := time.Now()
	var taskDate time.Time
	var err error

	if payload.Date == "" {
		taskDate = now
	} else {
		taskDate, err = time.Parse(dates.LayoutUser, payload.Date)
		if err != nil {
			taskDate, err = time.Parse(dates.LayoutDB, payload.Date)
			if err != nil {
				api.writeError(w, "неверный формат даты, ожидается DD.MM.YYYY или YYYYMMDD", http.StatusBadRequest)
				return
			}
		}
	}
	payload.Date = taskDate.Format(dates.LayoutDB)

	if payload.Repeat != "" {
		nextDateStr, err := NextDate(now, payload.Date, payload.Repeat)
		if err != nil {
			api.writeError(w, "некорректное правило повторения: "+err.Error(), http.StatusBadRequest)
			return
		}
		if taskDate.Before(now.Truncate(24 * time.Hour)) {
			payload.Date = nextDateStr
		}
	} else {
		if taskDate.Before(now.Truncate(24 * time.Hour)) {
			payload.Date = now.Format(dates.LayoutDB)
		}
	}

	taskToSave := db.Task{
		Date:    payload.Date,
		Title:   payload.Title,
		Comment: payload.Comment,
		Repeat:  payload.Repeat,
	}

	id, err := api.storage.AddTask(taskToSave)
	if err != nil {
		api.writeError(w, "не удалось добавить задачу в БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	api.writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)}, http.StatusCreated)
}
