package api

import (
	"go-final-project/pkg/db"
	"net/http"
)

// tasksResponse — это структура для ответа JSON со списком задач.
type tasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запросы на получение списка задач.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр поиска из URL
	search := r.URL.Query().Get("search")

	// Получаем задачи из БД с учетом поиска, ограничивая выборку 50 записями.
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeError(w, "не удалось получить задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ со списком задач.
	writeJSON(w, tasksResponse{Tasks: tasks}, http.StatusOK)
}
