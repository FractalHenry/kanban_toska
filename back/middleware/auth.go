package middleware

import (
	"backend/utils"
	"net/http"
	"strings"
)

// AuthMiddleware проверяет наличие и валидность JWT в запросах к защищенным маршрутам
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Не предоставлен токен", http.StatusUnauthorized)
			return
		}
		// Извлекаем токен из заголовка
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		// Проверяем валидность токена
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Не верный токен", http.StatusUnauthorized)
			return
		}
		// Добавляем логин пользователя в заголовок запроса
		r.Header.Set("login", claims.Login)
		// Передаем управление следующему обработчику в цепочке
		next.ServeHTTP(w, r)
	})
}
