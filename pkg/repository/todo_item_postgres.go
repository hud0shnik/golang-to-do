package repository

import (
	"fmt"
	"strings"
	"todo-app"

	"github.com/jmoiron/sqlx"
)

// Структура репозитория
type TodoItemPostgres struct {
	db *sqlx.DB
}

// Конструктор
func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

// Метод создания пункта
func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {

	// Создание транзакции
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// Id созданного пункта
	var itemId int

	// Запрос для создания записи в таблицу todo_items
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	// Выполнение запроса
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Запрос для вставки айди списка и пункта в таблицу lists_items
	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)

	// Выполнение запроса
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Приминяет изменения к БД и возвращает id пункта
	return itemId, tx.Commit()
}

// Метод получения всех пунктов по айди списка
func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {

	// Слайс пунктов
	var items []todo.TodoItem

	// Запрос на получение выборки пунктов
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id = li.list_id WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	// Запись найденных пунктов
	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	// Возвращает пункты и nil
	return items, nil
}

// Метод получения пункта по айди
func (r *TodoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {

	// Пункт, который будет возвращать метод
	var item todo.TodoItem

	// Запрос на получение пункта по айди
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	// Запись найденного пункта
	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	// Возвращает пункт и nil
	return item, nil
}

// Метод удаления пункта
func (r *TodoItemPostgres) Delete(userId, itemId int) error {

	// Запрос на удаление пункта с конкретным id
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	// Выполнение запроса
	_, err := r.db.Exec(query, userId, itemId)

	// Возвращает ошибку
	return err
}

// Метод обновления пункта
func (r *TodoItemPostgres) Update(userId, itemId int, input todo.UpdateItemInput) error {

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

	// Если поле готовности не пустое
	if input.Done != nil {

		// Добавляет в слайс элемент для формирования запроса к БД
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))

		// Добавляет в слайс аргументов значение нового описания
		args = append(args, *input.Done)

		// Инкремент id аргумента
		argId++
	}

	// Переменная для записи запроса к БД
	setQuery := strings.Join(setValues, ", ")

	// Формирование запроса
	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	// Исполнение запроса
	_, err := r.db.Exec(query, args...)

	// Возвращает ошибку
	return err
}
