package service

import (
	"todo-app"
	"todo-app/pkg/repository"
)

// Сервис для работы с пунктами
type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

// Конструктор
func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

// Метод передачи данных в слой repository при создании пункта
func (s *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {

	// Проверка на существование списка и доступ к нему
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	// Возвращает айди пункта и ошибку
	return s.repo.Create(listId, item)
}

// Метод передачи данных в слой repository при получении всех пунктов
func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

// Метод передачи данных в слой repository при получении пункта по айди
func (s *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

// Метод передачи данных в слой repository при удалении пункта
func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

// Метод передачи данных в слой repository при обновлении пункта
func (s *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
