package repository

import (
	"fmt"
	"strings"
	"todo-app"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Структура репозитория
type TodoListPostgres struct {
	db *sqlx.DB
}

// Конструктор
func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

// Метод создания списка
func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {

	// Создание транзакции
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// Id созданного списка
	var id int

	// Запрос для создания записи в таблицу todo_lists
	createListQuery := fmt.Sprintf("INSERT INTO %s (title,description) VALUES ($1,$2) RETURNING id", todoListTable)

	// Выполнение запроса
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Запрос для вставки id пользователя и id нового списка в таблицу users_lists
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id,list_id) VALUES ($1,$2)", usersListsTable)

	// Выполнение запроса
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Приминяет изменения к БД и возвращает id списка
	return id, tx.Commit()
}

// Метод получения всех списков по id пользователя
func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {

	// Слайс списков
	var lists []todo.TodoList

	// Запрос на получение выборки списков
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListTable, usersListsTable)

	// Запись найденных списков
	err := r.db.Select(&lists, query, userId)

	// Возвращает списки и ошибку
	return lists, err
}

// Метод получения списка по id
func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {

	// Список, который будет возвращать метод
	var list todo.TodoList

	// Запрос на получение списка по id
	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListTable, usersListsTable)

	// Запись найденного списка
	err := r.db.Get(&list, query, userId, listId)

	// Возвращает список и ошибку
	return list, err
}

// Метод удаления списка
func (r *TodoListPostgres) Delete(userId, listId int) error {

	// Запрос на удаление списка с конкретным id
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		todoListTable, usersListsTable)

	// Выполнение запроса
	_, err := r.db.Exec(query, userId, listId)

	// Возвращает ошибку
	return err
}

// Метод изменения списка
func (r *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {

	// Слайс значений
	setValues := make([]string, 0)

	// Слайс аргументов
	args := make([]interface{}, 0)

	// Id аргумента
	argId := 1

	// Если поле названия не пустое
	if input.Title != nil {

		// Добавляет в слайс элемент для формирования запроса к БД
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))

		// Добавляет в слайс аргументов значение нового названия
		args = append(args, *input.Title)

		// Инкремент id аргумента
		argId++
	}

	// Если поле описания не пустое
	if input.Description != nil {

		// Добавляет в слайс элемент для формирования запроса к БД
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))

		// Добавляет в слайс аргументов значение нового описания
		args = append(args, *input.Description)

		// Инкремент id аргумента
		argId++
	}

	// Переменная для записи запроса к БД
	setQuery := strings.Join(setValues, ", ")

	// Формирование запроса
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	// Исполнение запроса
	_, err := r.db.Exec(query, args...)

	return err
}
