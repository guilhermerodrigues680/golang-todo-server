package inmem

import (
	"fmt"
	"todoapp"
	"todoapp/storage"

	"github.com/sirupsen/logrus"
)

type StorageTodoInmem struct {
	todos         []todoapp.Todo
	currentSerial int
	logger        *logrus.Entry
}

func NewStorageTodoInmem(logger *logrus.Entry) *StorageTodoInmem {
	return &StorageTodoInmem{
		todos:         make([]todoapp.Todo, 0),
		currentSerial: 0,
		logger:        logger,
	}
}

func (stim *StorageTodoInmem) Create(description string) (todoapp.Todo, error) {
	stim.currentSerial++
	newTodo := todoapp.Todo{
		ID:          stim.currentSerial,
		Description: description,
		Done:        false,
	}
	stim.todos = append(stim.todos, newTodo)
	return newTodo, nil
}

func (stim *StorageTodoInmem) Read(id int) (todoapp.Todo, error) {
	for _, todo := range stim.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return todoapp.Todo{}, fmt.Errorf("%w : No item with id '%d'", storage.ErrNotFound, id)
}

func (stim *StorageTodoInmem) ReadAll() ([]todoapp.Todo, error) {
	return stim.todos, nil
}

func (stim *StorageTodoInmem) Update(id int, description string, done bool) (todoapp.Todo, error) {
	for idx, todo := range stim.todos {
		if todo.ID == id {
			updatedTodo := todoapp.Todo{
				ID:          id,
				Description: description,
				Done:        done,
			}
			stim.todos = append(stim.todos[:idx], stim.todos[idx+1:]...)
			stim.todos = append(stim.todos, updatedTodo)
			return updatedTodo, nil
		}
	}
	return todoapp.Todo{}, fmt.Errorf("%w : No item with id '%d'", storage.ErrNotFound, id)
}

func (stim *StorageTodoInmem) Delete(id int) error {
	for idx, todo := range stim.todos {
		if todo.ID == id {
			stim.todos = append(stim.todos[:idx], stim.todos[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w : No item with id '%d'", storage.ErrNotFound, id)
}
