package api

import (
	"net/http"
	"time"
)

// Init регистрирует все обработчики API на предоставленном маршрутизаторе.
func Init(mux *http.ServeMux) {
	// Открытые маршруты, доступные всем
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/signin", signInHandler)

	// Защищенные маршруты, обернутые в authMiddleware
	mux.Handle("/api/task", authMiddleware(http.HandlerFunc(taskHandler)))
	mux.Handle("/api/tasks", authMiddleware(http.HandlerFunc(tasksHandler)))
	mux.Handle("/api/task/done", authMiddleware(http.HandlerFunc(doneTaskHandler)))
}

// taskHandler - это мультиплексор для всех запросов /api/task.
// Он вызывает нужный обработчик в зависимости от HTTP-метода.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// nextDateHandler обрабатывает запросы на вычисление следующей даты.
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "неверный формат параметра 'now'", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
