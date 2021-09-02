package repository

type Autorization interface {
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Autorization
	TodoList
	TodoItem
}

func NewRepository() *Repository {
	return &Repository{}
}
