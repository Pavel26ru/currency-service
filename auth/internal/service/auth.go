package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/bcrypt"
)

var FailedAuthCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "Authorizations_Failed",
		Help: "Total number of failed authorizations",
	},
)

func init() {
	prometheus.MustRegister(FailedAuthCounter)
}

// UserService — структура для работы с пользователями
type UserService struct {
	jwtSecret []byte
	db        *sql.DB
}

// NewUserService — конструктор
func NewUserService(secret string, db *sql.DB) *UserService {
	return &UserService{
		jwtSecret: []byte(secret),
		db:        db,
	}
}

// Authenticate — проверка логина и пароля
func (s *UserService) Authenticate(login, password string) bool {
	var exists bool
	err := s.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)",
		login,
	).Scan(&exists)

	if err != nil || !exists {
		FailedAuthCounter.Inc()
		return false
	}

	var passwordHash string
	err = s.db.QueryRow(
		"SELECT password_hash FROM users WHERE login = $1",
		login,
	).Scan(&passwordHash)

	if err != nil {
		FailedAuthCounter.Inc()
		return false
	}

	success := checkPassword(password, passwordHash)
	if !success {
		FailedAuthCounter.Inc()
	}

	return success
}

// проверяем пароль
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
