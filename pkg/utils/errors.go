package utils

import "github.com/lib/pq"

// errors
var (
	ErrTimeout           = "timeout error ❌"
	ErrDBMigrations      = "failed to run migrations"
	ErrApiInitial        = "Api initial error"
	ErrIncorrectPassword = "incorrect password"
)

// warnings
var (
	WarnDBNotConnected = "database is not connected"
)

var (
	PqDuplicateErrorCode pq.ErrorCode = "23505"
	PqNoRowsFound                     = "sql: no rows in result set"
)
