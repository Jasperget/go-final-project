package api

import (
	"encoding/json"
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

// addTaskHandler обрабатывает запрос на добавление новой задачи.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var payload TaskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "ошибка десериализации JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 1. Валидация заголовка
	if payload.Title == "" {
		writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// 2. Обработка и валидация даты
	now := time.Now()
	var taskDate time.Time
	var err error

	if payload.Date == "" {
		taskDate = now
	} else {
		// Сначала пытаемся разобрать формат от браузера (DD.MM.YYYY)
		taskDate, err = time.Parse("02.01.2006", payload.Date)
		if err != nil {
			// Если не вышло, пытаемся разобрать формат для тестов (YYYYMMDD)
			taskDate, err = time.Parse(DateFormat, payload.Date)
			if err != nil {
				writeError(w, "неверный формат даты, ожидается DD.MM.YYYY или YYYYMMDD", http.StatusBadRequest)
				return
			}
		}
	}
	// Преобразуем дату в формат для хранения в БД (YYYYMMDD)
	payload.Date = taskDate.Format(DateFormat)

	// 3. Обработка правила повторения и даты
	if payload.Repeat != "" {
		// Если есть правило повторения, валидируем его и вычисляем следующую дату
		nextDateStr, err := NextDate(now, payload.Date, payload.Repeat)
		if err != nil {
			// Если NextDate вернула ошибку, значит правило некорректно
			writeError(w, "некорректное правило повторения: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Если исходная дата в прошлом, используем вычисленную следующую дату
		if taskDate.Before(now.Truncate(24 * time.Hour)) {
			payload.Date = nextDateStr
		}
	} else {
		// Если правила повторения нет, а дата в прошлом, используем сегодняшнюю дату
		if taskDate.Before(now.Truncate(24 * time.Hour)) {
			payload.Date = now.Format(DateFormat)
		}
	}

	// 4. Создание и сохранение задачи
	taskToSave := db.Task{
		Date:    payload.Date,
		Title:   payload.Title,
		Comment: payload.Comment,
		Repeat:  payload.Repeat,
	}

	id, err := db.AddTask(taskToSave)
	if err != nil {
		writeError(w, "не удалось добавить задачу в БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Отправка успешного ответа
	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)}, http.StatusCreated)
}

// writeJSON сериализует данные в JSON и отправляет клиенту.
func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError формирует и отправляет JSON с ошибкой.
func writeError(w http.ResponseWriter, message string, statusCode int) {
	writeJSON(w, map[string]string{"error": message}, statusCode)
}
