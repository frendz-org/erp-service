package postgres

import (
	stderrors "errors"
	"strings"

	"iam-service/pkg/errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func translateError(err error, entityName string) error {
	if err == nil {
		return nil
	}

	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return errors.ErrNotFound(entityName + " not found")
	}

	if stderrors.Is(err, gorm.ErrDuplicatedKey) {
		return errors.ErrConflict(entityName + " already exists")
	}

	var pgErr *pgconn.PgError
	if stderrors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return errors.ErrConflict(entityName + " already exists")
		case "23503":
			return errors.ErrConflict("referenced record does not exist")
		case "40001", "40P01":
			return errors.ErrConflict("database conflict, please retry")
		}
	}

	if isConnectionError(err) {
		return errors.ErrInternal("database connection error").WithError(err)
	}

	return errors.ErrInternal("database operation failed").WithError(err)
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	patterns := []string{
		"connection refused",
		"connection reset",
		"broken pipe",
		"no connection",
	}
	for _, p := range patterns {
		if strings.Contains(msg, p) {
			return true
		}
	}
	return false
}
