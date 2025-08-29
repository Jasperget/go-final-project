package api

import (
	"go-final-project/pkg/db"
	"net/http"
)

// tasksRequestLimit определяет максимальное количество задач, возвращаемых в одном запросе.
const tasksRequestLimit = 50

// tasksResponse — это структура для ответа JSON со списком задач.
type tasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

// tasksHandler теперь является методом *API.
func (api *API) tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := api.storage.Tasks(tasksRequestLimit, search)
	if err != nil {
		api.writeError(w, "не удалось получить задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	api.writeJSON(w, tasksResponse{Tasks: tasks}, http.StatusOK)
}
