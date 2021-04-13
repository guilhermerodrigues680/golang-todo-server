package postgres

import (
	"context"
	"todoapp"
	"todoapp/storage/postgres/query"

	"github.com/sirupsen/logrus"
)

type PostgresTodo struct {
	ps        *PostgresService
	logger    *logrus.Entry
	todoQuery *query.TodoQuery
}

func NewPostgresTodo(ps *PostgresService, logger *logrus.Entry) *PostgresTodo {
	return &PostgresTodo{
		ps:        ps,
		logger:    logger,
		todoQuery: query.NewTodoQuery(),
	}
}

func (pt *PostgresTodo) Read(id int) (todoapp.Todo, error) {
	sql, args, err := pt.todoQuery.SelectTodoById(id)
	if err != nil {
		return todoapp.Todo{}, err
	}

	todo := todoapp.Todo{}

	err = pt.ps.ConnPool.QueryRow(context.Background(), sql, args...).
		Scan(&todo.ID, &todo.Description, &todo.Done)
	if err != nil {
		pt.logger.Error(err)
		return todoapp.Todo{}, err
	}

	return todo, nil
}

func (pt *PostgresTodo) Create(description string) (todoapp.Todo, error) {
	sql, args, err := pt.todoQuery.InsertTodo(description, false)
	if err != nil {
		return todoapp.Todo{}, err
	}

	var id int

	err = pt.ps.ConnPool.QueryRow(context.Background(), sql, args...).
		Scan(&id)
	if err != nil {
		pt.logger.Error(err)
		return todoapp.Todo{}, err
	}

	todo, err := pt.Read(id)
	if err != nil {
		return todoapp.Todo{}, err
	}

	return todo, nil
}

func (pt *PostgresTodo) ReadAll() ([]todoapp.Todo, error) {
	sql, args, err := pt.todoQuery.SelectAllTodo()
	if err != nil {
		return nil, err
	}

	rows, err := pt.ps.ConnPool.Query(context.Background(), sql, args...)

	todoList := make([]todoapp.Todo, 0)
	for rows.Next() {
		todo := todoapp.Todo{}
		if err := rows.Scan(&todo.ID, &todo.Description, &todo.Done); err != nil {
			pt.logger.Errorf("Row scan failed: %v", err)
			return nil, err
		}
		todoList = append(todoList, todo)
	}

	return todoList, nil
}

func (pt *PostgresTodo) Update(id int, description string, done bool) (todoapp.Todo, error) {
	sql, args, err := pt.todoQuery.UpdateTodoById(id, description, done)
	if err != nil {
		return todoapp.Todo{}, err
	}

	commandTag, err := pt.ps.ConnPool.Exec(context.Background(), sql, args...)
	if err != nil {
		pt.logger.Errorf("%s, %v", commandTag, err)
		return todoapp.Todo{}, err
	}

	todo, err := pt.Read(id)
	if err != nil {
		return todoapp.Todo{}, err
	}

	return todo, nil
}

func (pt *PostgresTodo) Delete(id int) error {
	sql, args, err := pt.todoQuery.DeleteTodoById(id)
	if err != nil {
		return err
	}

	commandTag, err := pt.ps.ConnPool.Exec(context.Background(), sql, args...)
	if err != nil {
		pt.logger.Errorf("%s, %v", commandTag, err)
		return err
	}

	return nil
}
