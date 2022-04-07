package repository

import (
	"fmt"
	"todo-app"

	"github.com/jmoiron/sqlx"
)

// Структура, которая имплементирует интерфейс репозитория и работает с БД
type AuthPostgres struct {
	db *sqlx.DB
}

// Конструктор
func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// Метод создания пользователя
func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {

	// Id, который будет возвращать метод
	var id int

	// Создание запроса для добавления пользователя в БД
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", userTable)

	// Выполнение запроса
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)

	// Обработка ошибки
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	// Возвращает айди и nil
	return id, nil
}

// Метод получения пользователя
func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {

	// Пользователь, которого будет возвращать метод
	var user todo.User

	// Создание запроса для поиска пользователя
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", userTable)

	// Выполнение запроса
	err := r.db.Get(&user, query, username, password)

	// Возвращает пользователя и nil
	return user, err
}
