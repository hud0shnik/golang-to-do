package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Константы названий таблиц
const (
	userTable       = "users"
	todoListTable   = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

// Структура с параметрами для подключения к БД
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// Функция инициализации базы данных
func NewPostgresDB(cfg Config) (*sqlx.DB, error) {

	// Создание БД
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	// Проверка подключения к БД
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
