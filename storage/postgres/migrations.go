package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

// Migrations é uma estrutura utilitária para migrações de tabelas
type Migrations struct {
	connPool *pgxpool.Pool
	logger   *logrus.Entry
}

func NewMigrations(connPool *pgxpool.Pool, logger *logrus.Entry) *Migrations {
	return &Migrations{
		connPool: connPool,
		logger:   logger,
	}
}

func (m *Migrations) Start() error {
	err := m.createTodoTable(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (m *Migrations) createTodoTable(ctx context.Context) error {
	sql := `CREATE TABLE TODO (
		TODO_ID SERIAL PRIMARY KEY,
		TODO_DESCRIPTION VARCHAR(200) NOT NULL,
		TODO_DONE BOOL NOT NULL,
		TODO_UPDATE_AT TIMESTAMP NOT NULL DEFAULT (NOW()),
		TODO_CREATED_AT TIMESTAMP NOT NULL DEFAULT (NOW())
	)`

	commandTag, err := m.connPool.Exec(ctx, sql)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); IsDuplicateTableError(pgErr) {
			m.logger.Warnf("%s, migration was skipped", pgErr.Message)
			return nil
		}

		m.logger.Error(err)
		return err
	}

	m.logger.Infof("%s, successfully created table TODO", commandTag)

	return nil
}
