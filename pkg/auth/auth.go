package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	tokenTTL   = 8 * time.Hour
	signingKey = "a_very_secret_key_that_should_be_long_and_random" // В реальном проекте этот ключ должен быть сложнее и храниться в секрете
)

// claims — это структура полезной нагрузки нашего токена.
type claims struct {
	jwt.RegisteredClaims
	PasswordHash string `json:"password_hash"`
}

// GenerateToken создает новый JWT-токен на основе пароля.
func GenerateToken(password string) (string, error) {
	// Создаем хэш пароля, как требует задание
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// Устанавливаем время жизни токена
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		PasswordHash: hash,
	})

	return token.SignedString([]byte(signingKey))
}

// ValidateToken проверяет токен и сверяет хэш пароля.
func ValidateToken(tokenString, currentPassword string) (bool, error) {
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Убеждаемся, что алгоритм подписи тот, который мы ожидаем
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, fmt.Errorf("token is not valid")
	}

	// Сравниваем хэш из токена с хэшем ТЕКУЩЕГО пароля из переменной окружения
	currentHash := fmt.Sprintf("%x", sha256.Sum256([]byte(currentPassword)))
	if claims.PasswordHash != currentHash {
		return false, fmt.Errorf("password has changed, token is invalid")
	}

	return true, nil
}
