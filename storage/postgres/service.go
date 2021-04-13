package postgres

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresService struct {
	ConnPool   *pgxpool.Pool
	connString string
	logger     *logrus.Entry
}

type pgsqlLogger struct {
	*logrus.Entry
}

func NewPgsqlLogger(logger *logrus.Entry) *pgsqlLogger {
	return &pgsqlLogger{
		Entry: logger,
	}
}

func (m *pgsqlLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	m.Trace(msg)
	if sql, ok := data["sql"]; ok {
		m.Tracef("Executing SQL: %s", sql)
	}
}

// NewPostgresService retorna um cliente capaz de se comunicar com o Banco de Dados
func NewPostgresService(host string, port int, dbname string, user string, password string, logger *logrus.Entry) (*PostgresService, error) {
	connString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", host, port, dbname, user, password)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Errorf("Unable to parse config: %v", err)
		return nil, err
	}

	internalLogger := logger.WithField("internal", "pgx")
	config.ConnConfig.Logger = NewPgsqlLogger(internalLogger)
	config.ConnConfig.LogLevel = pgx.LogLevelTrace
	// config.LazyConnect = true
	// config.MinConns = 4
	// config.MaxConns = 4

	config.BeforeConnect = func(c context.Context, cc *pgx.ConnConfig) error {
		internalLogger.Tracef("BeforeConnect")
		return nil
	}
	config.AfterConnect = func(c1 context.Context, c2 *pgx.Conn) error {
		internalLogger.Tracef("AfterConnect")
		return nil
	}
	config.BeforeAcquire = func(c1 context.Context, c2 *pgx.Conn) bool {
		internalLogger.Tracef("BeforeAcquire")
		return true
	}
	config.AfterRelease = func(c *pgx.Conn) bool {
		internalLogger.Tracef("AfterRelease")
		return true
	}

	// dbpool, err := pgxpool.Connect(context.Background(), connString)
	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		logger.Errorf("Unable to connect to database: %v", err)
		return nil, err
	}
	// defer dbpool.Close()
	logger.Trace("Database opened successfully")

	err = dbpool.Ping(context.Background())
	if err != nil {
		logger.Errorf("Unable to ping: %v", err)
		return nil, err
	}
	logger.Trace("database successfully pinged")
	logger.Info("Successfully connected to the database")

	// var greeting string
	// err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(greeting)

	return &PostgresService{
		ConnPool:   dbpool,
		connString: connString,
		logger:     logger,
	}, nil
}

func (ps *PostgresService) MigrateTables() error {
	m := NewMigrations(ps.ConnPool, ps.logger.WithField("internal", "migrations"))
	err := m.Start()
	if err != nil {
		return err
	}

	return nil
}
