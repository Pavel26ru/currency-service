package repository

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("user not found")
)

type User struct {
	Login    string
	Password string
}

type UserRepository struct {
	db *sql.DB
	mu *sync.RWMutex
}

func NewUser(db *sql.DB) UserRepository {
	return UserRepository{
		db: db,
		mu: &sync.RWMutex{},
	}
}

// хешируем пароль с солью
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// проверяем пароль
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (repo *UserRepository) AddUser(user User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var exists bool
	// проверяем, существует ли такой пользователь
	err := repo.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)", user.Login).Scan(&exists)

	if err != nil {
		return err
	}
	if exists {
		return ErrUserAlreadyExist
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	// добавляем пользователя в БД
	_, err = repo.db.Exec(
		"INSERT INTO users (login, password_hash) VALUES ($1, $2)",
		user.Login, hashedPassword,
	)

	return err
}

func (repo *UserRepository) GetUser(ctx context.Context, login string) (User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	var user User
	err := repo.db.QueryRowContext(
		ctx,
		"SELECT id, login, password_hash FROM users WHERE login = $1",
		login,
	).Scan(&user.Login, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

// проверяет логин и пароль
func (repo *UserRepository) AuthenticateUser(ctx context.Context, login, password string) (bool, error) {
	user, err := repo.GetUser(ctx, login)
	if err != nil {
		if err == ErrUserNotFound {
			return false, nil
		}
		return false, err
	}

	return checkPassword(password, user.Password), nil
}
