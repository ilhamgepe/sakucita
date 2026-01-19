package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func PgConstraint(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.ConstraintName
	}
	return ""
}

// Postgres error codes reference:
// https://www.postgresql.org/docs/current/errcodes-appendix.html

const (
	errCodeDuplicateUniqueViolation = "23505"
	errCodeForeignKeyViolation      = "23503"
	errCodeNotNullViolation         = "23502"
	errCodeExclusionViolation       = "23P01"
)

// IsPgErrCode checks if the given error is a Postgres error with the specified code.
func IsPgErrCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == code
	}
	return false
}

// IsUniqueViolation checks if the error is a unique constraint violation.
func IsDuplicateUniqueViolation(err error) bool {
	return IsPgErrCode(err, errCodeDuplicateUniqueViolation)
}

// IsForeignKeyViolation checks if the error is a foreign key constraint violation.
func IsForeignKeyViolation(err error) bool {
	return IsPgErrCode(err, errCodeForeignKeyViolation)
}
