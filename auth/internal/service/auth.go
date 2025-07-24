package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserService — структура для работы с пользователями
type UserService struct {
	// Здесь могут быть зависимости: БД, конфиг, секрет для JWT и т.д.
	jwtSecret []byte
	// users map[string]string // если хранишь пользователей в памяти: login -> password
}

// NewUserService — конструктор
func NewUserService(secret string) *UserService {
	return &UserService{
		jwtSecret: []byte(secret),
		// users: map[string]string{"admin": "admin"}, // пример
	}
}

// Authenticate — проверка логина и пароля
func (s *UserService) Authenticate(login, password string) bool {
	// TODO: Реализуй свою логику (например, запрос к БД)
	// Пример для in-memory:
	// realPassword, ok := s.users[login]
	// return ok && realPassword == password
	return login == "admin" && password == "admin" // временная заглушка
}

// GenerateJWT — генерация JWT-токена
func (s *UserService) GenerateJWT(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(s.jwtSecret)
}

// ValidateJWT — валидация токена (если потребуется)
func (s *UserService) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	login, ok := claims["login"].(string)
	if !ok {
		return "", errors.New("login not found in token")
	}
	return login, nil
}
