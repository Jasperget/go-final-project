package api

import (
	"encoding/json"
	"go-final-project/pkg/auth"
	"go-final-project/pkg/config"
	"net/http"
)

// signInHandler теперь является методом *API.
func (api *API) signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	storedPassword := config.Get().Password
	if storedPassword == "" {
		api.writeError(w, "аутентификация не включена на сервере", http.StatusInternalServerError)
		return
	}

	var payload struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		api.writeError(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if payload.Password != storedPassword {
		api.writeError(w, "неверный пароль", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(payload.Password)
	if err != nil {
		api.writeError(w, "не удалось создать токен", http.StatusInternalServerError)
		return
	}

	api.writeJSON(w, map[string]string{"token": token}, http.StatusOK)
}

// authMiddleware теперь является методом *API.
func (api *API) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		password := config.Get().Password
		if password == "" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "требуется аутентификация", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		valid, err := auth.ValidateToken(tokenStr, password)
		if !valid || err != nil {
			http.Error(w, "невалидный токен", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
