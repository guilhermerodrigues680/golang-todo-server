package todoapp

import (
	"errors"
	"fmt"
	"todoapp/storage"

	"github.com/sirupsen/logrus"
)

// Representa um TO-DO no sistema
type Todo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type TodoStorage interface {
	Create(description string) (Todo, error)
	Read(id int) (Todo, error)
	ReadAll() ([]Todo, error)
	Update(id int, description string, done bool) (Todo, error)
	Delete(id int) error
}

type TodoService struct {
	s      TodoStorage
	logger *logrus.Entry
}

// ErrNotFound é o erro returnado pelo service quando uma busca no storage não retorna resultados.
var ErrNotFound = storage.ErrNotFound

func NewTodoService(s TodoStorage, logger *logrus.Entry) *TodoService {
	return &TodoService{
		s:      s,
		logger: logger,
	}
}

func (ts *TodoService) Create(description string) (Todo, error) {
	return ts.s.Create(description)
}

func (ts *TodoService) Read(id int) (Todo, error) {
	todo, err := ts.s.Read(id)
	if errors.Is(err, storage.ErrNotFound) {
		return todo, fmt.Errorf("%w", ErrNotFound)
	} else if err != nil {
		return todo, err
	}

	return todo, nil
}

func (ts *TodoService) ReadAll() ([]Todo, error) {
	return ts.s.ReadAll()
}

func (ts *TodoService) Update(id int, description string, done bool) (Todo, error) {
	return ts.s.Update(id, description, done)
}

func (ts *TodoService) Delete(id int) error {
	return ts.s.Delete(id)
}
