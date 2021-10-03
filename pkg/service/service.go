package service

import (
	"todo-app"
	"todo-app/pkg/repository"
)

type Autorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Autorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Autorization: NewAuthService(repos.Autorization),
	}
}
