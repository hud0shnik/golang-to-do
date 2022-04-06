// Файл со структурами и функциями, которые используются во всём приложении
// Находится в корне проекта для удобства

package todo

import "errors"

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}

// Структура ввода при обновлении списка
type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

// Проверака при обновлении списка
func (i UpdateListInput) Validate() error {

	// Если оба значения пустые, возвращает ошибку
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	// Возвращает nil
	return nil
}

// Структура ввода при обновлении пункта
type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

// Проверака при обновлении пункта
func (i UpdateItemInput) Validate() error {

	// Если оба значения пустые, возвращает ошибку
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}

	// Возвращает nil
	return nil
}
