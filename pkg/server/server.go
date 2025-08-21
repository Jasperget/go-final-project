package server

import (
	"go-final-project/pkg/api"
	"net/http"
)

// Run теперь принимает обработчик API как аргумент.
func Run(port string, apiHandler *api.API) error {
	mux := http.NewServeMux()

	// Регистрируем обработчики API
	apiHandler.RegisterRoutes(mux)

	// Настраиваем обработку статических файлов
	fileServer := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fileServer)

	return http.ListenAndServe(":"+port, mux)
}
