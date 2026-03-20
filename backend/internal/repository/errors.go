package repository

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// IsDuplicateEntry devuelve true si el error es un MySQL 1062 (Duplicate entry).
// Usar esto en vez de strings.Contains para mayor robustez ante distintas versiones del driver.
func IsDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
