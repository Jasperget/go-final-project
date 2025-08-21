package api

import (
	"encoding/json"
	"go-final-project/pkg/auth"
	"go-final-project/pkg/config" // Импортируем пакет config
	"net/http"
)

// signInHandler обрабатывает POST-запрос на аутентификацию.
func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Используем пароль из загруженной конфигурации
	storedPassword := config.Get().Password
	if storedPassword == "" {
		writeError(w, "аутентификация не включена на сервере", http.StatusInternalServerError)
		return
	}

	var payload struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if payload.Password != storedPassword {
		writeError(w, "неверный пароль", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(payload.Password)
	if err != nil {
		writeError(w, "не удалось создать токен", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": token}, http.StatusOK)
}

// authMiddleware проверяет аутентификацию для защищенных маршрутов.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Используем пароль из загруженной конфигурации
		password := config.Get().Password
		if password == "" {
			// Если пароль не установлен, аутентификация не требуется, пропускаем дальше
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			// Если куки нет, возвращаем ошибку 401
			http.Error(w, "требуется аутентификация", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		valid, err := auth.ValidateToken(tokenStr, password)
		if !valid || err != nil {
			// Если токен невалидный, возвращаем ошибку 401
			http.Error(w, "невалидный токен", http.StatusUnauthorized)
			return
		}

		// Если все в порядке, пропускаем запрос к основному обработчику
		next.ServeHTTP(w, r)
	})
}
