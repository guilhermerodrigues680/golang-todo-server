package postgres

import "github.com/jackc/pgconn"

func IsDuplicateTableError(err *pgconn.PgError) bool {
	return err.Code == "42P07"
}
