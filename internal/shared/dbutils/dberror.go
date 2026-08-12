package dbutils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const UniqueViolationCode = "23505"

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == UniqueViolationCode
}

func ConstraintName(err error) string {
	if err == nil {
		return ""
	}

	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return ""
	}

	return pgErr.ConstraintName
}
