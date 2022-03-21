package service

import (
	"todo-app"
	"todo-app/pkg/repository"
)

// Сервис для работы со списками
type TodoListService struct {
	repo repository.TodoList
}

// Конструктор
func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

// Метод передачи данных в слой repository при создании списка
func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

// Метод передачи данных в слой repository при получении списков
func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}

// Метод передачи данных в слой repository при получении списка по id
func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}

// Метод передачи данных в слой repository при удалении списка
func (s *TodoListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

// Метод передачи данных в слой repository при обновлении списка
func (s *TodoListService) Update(userId, listId int, input todo.UpdateListInput) error {

	// Валидация вводных данных перед передачей в repository
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}
