package query

import (
	sq "github.com/Masterminds/squirrel"
)

// TodoQuery é uma classe utilitária para construir consultas SQL
type TodoQuery struct {
	psql *sq.StatementBuilderType
}

func NewTodoQuery() *TodoQuery {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &TodoQuery{
		psql: &psql,
	}
}

func (tq *TodoQuery) selectTodo() sq.SelectBuilder {
	return tq.psql.
		Select("TODO_ID", "TODO_DESCRIPTION", "TODO_DONE").
		From("TODO")
}

func (tq *TodoQuery) SelectAllTodo() (string, []interface{}, error) {
	return tq.selectTodo().
		ToSql()
}

func (tq *TodoQuery) SelectTodoById(id int) (string, []interface{}, error) {
	return tq.selectTodo().
		Where(sq.Eq{"TODO_ID": id}).
		ToSql()
}

func (tq *TodoQuery) InsertTodo(description string, done bool) (string, []interface{}, error) {
	return tq.psql.
		Insert("TODO").
		Columns("TODO_DESCRIPTION", "TODO_DONE").
		Values(description, done).
		Suffix("RETURNING TODO_ID").
		ToSql()
}

func (tq *TodoQuery) UpdateTodoById(id int, description string, done bool) (string, []interface{}, error) {
	return tq.psql.Update("TODO").
		Set("TODO_DESCRIPTION", description).
		Set("TODO_DONE", done).
		Where(sq.Eq{"TODO_ID": id}).
		ToSql()
}

func (tq *TodoQuery) DeleteTodoById(id int) (string, []interface{}, error) {
	return tq.psql.Delete("TODO").
		Where(sq.Eq{"TODO_ID": id}).
		ToSql()
}
