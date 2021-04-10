package inmem

import (
	"fmt"
	"todoserver"
	"todoserver/storage"
)

type StorageTodoInmem struct {
	Todos         []todoserver.Todo
	CurrentSerial int
}

func NewStorageTodoInmem() *StorageTodoInmem {
	return &StorageTodoInmem{
		Todos:         make([]todoserver.Todo, 0),
		CurrentSerial: 0,
	}
}

func (stim *StorageTodoInmem) Create(description string) (todoserver.Todo, error) {
	stim.CurrentSerial++
	newTodo := todoserver.Todo{
		ID:          stim.CurrentSerial,
		Description: description,
		Done:        false,
	}
	stim.Todos = append(stim.Todos, newTodo)
	return newTodo, nil
}

func (stim *StorageTodoInmem) Read(id int) (todoserver.Todo, error) {
	for _, todo := range stim.Todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return todoserver.Todo{}, fmt.Errorf("%w '%d'", storage.ErrNotFound, id)
}

func (stim *StorageTodoInmem) ReadAll() ([]todoserver.Todo, error) {
	return stim.Todos, nil
}

func (stim *StorageTodoInmem) Update(id int, description string, done bool) (todoserver.Todo, error) {
	for idx, todo := range stim.Todos {
		if todo.ID == id {
			updatedTodo := todoserver.Todo{
				ID:          id,
				Description: description,
				Done:        done,
			}
			stim.Todos = append(stim.Todos[:idx], stim.Todos[idx+1:]...)
			stim.Todos = append(stim.Todos, updatedTodo)
			return updatedTodo, nil
		}
	}
	return todoserver.Todo{}, fmt.Errorf("%w '%d'", storage.ErrNotFound, id)
}

func (stim *StorageTodoInmem) Delete(id int) error {
	for idx, todo := range stim.Todos {
		if todo.ID == id {
			stim.Todos = append(stim.Todos[:idx], stim.Todos[idx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w '%d'", storage.ErrNotFound, id)
}
