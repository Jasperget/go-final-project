package api

import (
	"encoding/json"
	"go-final-project/pkg/dates"
	"go-final-project/pkg/db"
	"net/http"
	"time"
)

// API представляет собой обработчик API, который зависит от хранилища.
type API struct {
	storage *db.Storage
}

// New создает новый экземпляр API с предоставленным хранилищем.
func New(storage *db.Storage) *API {
	return &API{storage: storage}
}

// RegisterRoutes регистрирует все маршруты API на предоставленном маршрутизаторе.
// Эта функция заменяет старую Init.
func (api *API) RegisterRoutes(mux *http.ServeMux) {
	// Открытые маршруты, доступные всем
	mux.HandleFunc("/api/nextdate", api.nextDateHandler)
	mux.HandleFunc("/api/signin", api.signInHandler)

	// Защищенные маршруты, обернутые в authMiddleware
	mux.Handle("/api/task", api.authMiddleware(http.HandlerFunc(api.taskHandler)))
	mux.Handle("/api/tasks", api.authMiddleware(http.HandlerFunc(api.tasksHandler)))
	mux.Handle("/api/task/done", api.authMiddleware(http.HandlerFunc(api.doneTaskHandler)))
}

// taskHandler - это мультиплексор для всех запросов /api/task.
// Он вызывает нужный обработчик в зависимости от HTTP-метода.
func (api *API) taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		api.getTaskHandler(w, r)
	case http.MethodPost:
		api.addTaskHandler(w, r)
	case http.MethodPut:
		api.updateTaskHandler(w, r)
	case http.MethodDelete:
		api.deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// nextDateHandler обрабатывает запросы на вычисление следующей даты.
func (api *API) nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dates.LayoutDB, nowStr)
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

// writeJSON и writeError теперь методы *API, чтобы избежать дублирования.
func (api *API) writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (api *API) writeError(w http.ResponseWriter, message string, statusCode int) {
	api.writeJSON(w, map[string]string{"error": message}, statusCode)
}
